package framework

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/luopengift/framework/pkg/encoding/json"
	"github.com/luopengift/framework/util"
	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

// var default var
var (
	wg sync.WaitGroup
)

// App framework
type App struct {
	context.Context `json:"-"`
	*Option         `json:"option"`
	Name            string         `json:"name" yaml:"name"`
	ID              string         `json:"id" yaml:"id"`
	TimeZone        *time.Location `json:"timeZone"`
	config          interface{}
	onPrepare       Function
	onInit          Function
	onMain          Function
	onThreads       []Function
	goroutines      []*Goroutine
	onExit          Function
	errChan         chan error
}

// New new app instance
func New(opts ...*Option) *App {
	app := &App{
		Context: context.Background(),
		Option:  defaultOption,
		Name:    "",
		ID:      util.Random(10),
	}
	app.Option.Merge(opts...)
	app.MainFunc(DefaultMainThread)
	app.ExitFunc(defaultFunc)
	return app
}

func (app *App) init() {
	app.errChan = make(chan error, 100)
}

// WithContext with context
func (app *App) WithContext(ctx context.Context) {
	app.Context = ctx
}

func (app *App) Error(format string, v ...interface{}) {
	app.errChan <- fmt.Errorf(format, v...)
}

func (app *App) error() {
	for {
		select {
		case err, ok := <-app.errChan:
			if !ok {
				return
			}
			log.Error("%v", err)
		}
	}
}

// Bind bind runner interface
func (app *App) Bind(run Runner) {
	app.PrepareFunc(run.Prepare)
	app.InitFunc(run.Init)
	app.MainFunc(run.Main)
	app.ThreadFunc(run.Thread)
	app.ExitFunc(run.Exit)
}

// BindConfig bind config
func (app *App) BindConfig(v interface{}) {
	app.config = v
}

// Prepare Preparer interface
func (app *App) Prepare(prepare Function) {
	app.onPrepare = prepare
}

// PrepareFunc ParpareFunc
func (app *App) PrepareFunc(f Func) {
	app.onPrepare = f
}

// Init Initer interface
func (app *App) Init(init Function) {
	app.onInit = init
}

// InitFunc init func program global var
func (app *App) InitFunc(f Func) {
	app.onInit = f
}

// Main interface
func (app *App) Main(main Function) {
	app.onMain = main
}

// MainFunc handle main func
func (app *App) MainFunc(f Func) {
	app.onMain = f
}

// Exit Exiter interface
func (app *App) Exit(exit Function) {
	app.onExit = exit
}

// ExitFunc init func program global var
func (app *App) ExitFunc(f Func) {
	app.onExit = f
}

// Thread interface
func (app *App) Thread(threads ...Function) {
	for _, thread := range threads {
		app.onThreads = append(app.onThreads, thread)
	}
}

// ThreadFunc Func init func program global var
func (app *App) ThreadFunc(fs ...Func) {
	for _, f := range fs {
		app.onThreads = append(app.onThreads, f)
	}
}

// 用户管理goroutine运行退出等状态, framework仅调起函数
func (app *App) runThreads(ctx context.Context) error {
	for id, thread := range app.onThreads {
		if thread == nil {
			return log.Errorf("thread must not nil!")
		}
		wg.Add(1)
		go func(ctx context.Context, id int, execute Function) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					log.Fatal("Thread[%v] %v\n%v", id, err, string(debug.Stack()))
				}
			}()
			if err := execute.Func(ctx); err != nil {
				log.Error("Thread[%v]: %v", id, err)
			}
		}(ctx, id, thread)
	}
	return nil
}

// GoroutineFunc GoroutineFunc
func (app *App) GoroutineFunc(name string, fs FunctionWithExit, num ...int) {
	var min, max int
	switch len(num) {
	case 0:
		min, max = 1, 1
	case 1:
		min, max = num[0], num[0]
	case 2:
		min, max = num[0], num[1]
	default:
		min, max = -1, -1
	}
	if name == "" {
		name = util.Random(8)
	}
	app.goroutines = append(app.goroutines, newGoroutine(name, fs, min, max))
}

// 协程循环由framework管理
func (app *App) runGoroutines(ctx context.Context) error {
	for _, goroutine := range app.goroutines {
		if goroutine.exec == nil {
			return log.Errorf("goroutine must not nil!")
		}
		for i := 0; i < goroutine.min; i++ {
			wg.Add(1)
			go func(ctx context.Context, name string, seq, num int, gor *Goroutine) {
				defer wg.Done()
				var (
					exit  bool // true: 退出goroutine, false: 循环调用goroutine.
					err   error
					entry = func(seq int, execute FunctionWithExit) (bool, error) {
						defer func() {
							if err := recover(); err != nil {
								log.Fatal("goroutine panic[%v-%v/%v] %v\n%v", name, seq, num, err, string(debug.Stack()))
							}
						}()
						return execute.FuncWithExit(ctx)
					}
				)

				for !exit {
					select {
					case <-ctx.Done():
						log.Error("goroutine ctx[%v-%v/%v]: %v", name, seq, num, ctx.Err())
						return
					default:
						if err := gor.before.Func(ctx); err != nil {
							log.Error("before error: %v", err)
							break
						}
						if exit, err = entry(seq, gor.exec); err != nil {
							log.Error("goroutine run[%v-%v/%v]: %v", name, seq, num, err)
						}
						if err := gor.before.Func(ctx); err != nil {
							log.Error("after error: %v", err)
						}
					}
				}
			}(ctx, goroutine.name, i+1, goroutine.min, goroutine)
		}
	}
	return nil
}

// LoadConfig loading config step by step
func (app *App) LoadConfig() error {
	envOpt, err := newEnvOpt()
	if err != nil {
		return err
	}
	argsOpt := newArgsOpt()
	app.Option.Merge(envOpt, argsOpt) // 仅为了合并configPath供配置文件使用

	ok, err := util.PathExist(app.Option.ConfigPath)
	if err != nil {
		return err
	}
	if ok {
		if app.config != nil {
			if err := types.ParseConfigFile(app.config, app.Option.ConfigPath); err != nil {
				return err
			}
		}
		if err := types.ParseConfigFile(app.Option, app.Option.ConfigPath); err != nil {
			return err
		}
		if err := types.ParseConfigFile(app, app.Option.ConfigPath); err != nil {
			return err
		}
		app.Option.Merge(envOpt, argsOpt) // 修改被配置文件改掉配置
	}
	if app.config != nil {
		if err = json.Format(app.config, app.Option); err != nil {
			return err
		}
	}
	if app.Name == "" {
		app.Name = filepath.Base(os.Args[0])
	}
	app.TimeZone, err = time.LoadLocation(app.Option.Tz)
	return err
}

// Run app instance
func (app *App) Run() {
	now := time.Now()
	defer func(now time.Time) {
		if err := recover(); err != nil {
			log.Fatal("%v\n%v", err, string(debug.Stack()))
		}
		log.Warn("exit: running time=%v", time.Since(now))
	}(now)

	if err := app.execute(); err != nil {
		log.Error("%v", err)
	}
	if app.onExit != nil {
		if err := app.onExit.Func(app.Context); err != nil {
			log.Error("exit: %v", err)
		}
	}
	wg.Wait()
}

// 执行逻辑:
// 0. 声明NumCPU和context
// 1. 初始化app.init()
// 2. 调用onPrepare过程
// 3. 加载配置信息
// 4. 初始化日志模块, app.initLog()
// 5. 调用初始化onInit过程
// 6. 加载其他框架模块, 例如http, pprof等
// 7. 调用主函数onMain过程
// 8. 调用其他goroutines
// 9. 获取退出信号, sign.Kill或者mainExit
// 10. 调用退出onExit过程
func (app *App) execute() error {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx, cancel := context.WithCancel(app.Context)
	defer cancel()
	app.init()

	if app.onPrepare != nil {
		if err = app.onPrepare.Func(ctx); err != nil {
			return err
		}
	}

	if err = app.LoadConfig(); err != nil {
		return err
	}

	if app.Option.Version {
		log.ConsoleWithMagenta("%v", version.String())
		return nil
	}

	if err := app.InitLog(); err != nil {
		return err
	}

	log.Info("[%s] init...", app.Name)

	if app.onInit != nil {
		if err := app.onInit.Func(ctx); err != nil {
			return err
		}
	}

	// http
	app.initHttpd()

	// pprof
	if app.Option.PprofPath != "" {
		if err = os.MkdirAll(app.Option.PprofPath, 0755); err != nil {
			return err
		}
		cpu, err := os.Create(filepath.Join(app.Option.PprofPath, "cpu.prof"))
		if err != nil {
			return err
		}
		defer cpu.Close()

		if err = pprof.StartCPUProfile(cpu); err != nil {
			return err
		}
		defer pprof.StopCPUProfile()

		mem, err := os.Create(filepath.Join(app.Option.PprofPath, "mem.prof"))
		if err != nil {
			return err
		}
		defer mem.Close()
		//mem
		runtime.GC()
		if err := pprof.WriteHeapProfile(mem); err != nil {
			return err
		}
	}

	if app.onMain == nil {
		return log.Errorf("Main is nil, must set!")
	}
	log.Info("init success.")
	mainExit := make(chan struct{})
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("MainThread: %v\n%v", err, string(debug.Stack()))
			}
			mainExit <- struct{}{}
		}()
		if err := app.onMain.Func(ctx); err != nil {
			log.Error("MainThread: %v", err)
		}
	}(ctx)

	if err := app.runThreads(ctx); err != nil {
		return err
	}

	if err := app.runGoroutines(ctx); err != nil {
		return err
	}

	signSystem := make(chan os.Signal)
	signal.Notify(signSystem, os.Interrupt, os.Kill)

	select {
	case <-mainExit:
		log.Warn("mainThread exit...")
	case s := <-signSystem:
		log.Warn("Get signal: %v", s)
	}
	return nil
}
