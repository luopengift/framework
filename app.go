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
	"time"

	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

// var default var
var (
	TimeZone *time.Location
)

// App framework
type App struct {
	*Option
	Name          string `json:"name" yaml:"name"`
	ID            string `json:"id" yaml:"id"`
	config        interface{}
	onPrepare     Preparer
	onInit        Initer
	onMain        Mainer
	onThreads     []Threader
	onThreadLoops []Looper
	onExit        Exiter
	errChan       chan error
}

// New new app instance
func New(opts ...*Option) *App {
	app := &App{
		Option: defaultOption,
	}
	app.Option.Merge(opts...)
	return app
}

func (app *App) init() {
	app.errChan = make(chan error, 100)
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
	app.LoopFunc(run.Loop)
	app.ExitFunc(run.Exit)
}

// BindConfig bind config
func (app *App) BindConfig(v interface{}) {
	app.config = v
}

// Prepare Preparer interface
func (app *App) Prepare(prepare Preparer) {
	app.onPrepare = prepare
}

// PrepareFunc ParpareFunc
func (app *App) PrepareFunc(f PreparerFunc) {
	app.onPrepare = f
}

// Init Initer interface
func (app *App) Init(init Initer) {
	app.onInit = init
}

// InitFunc init func program global var
func (app *App) InitFunc(f IniterFunc) {
	app.onInit = f
}

// Main interface
func (app *App) Main(main Mainer) {
	app.onMain = main
}

// MainFunc handle main func
func (app *App) MainFunc(f MainerFunc) {
	app.onMain = f
}

// Exit Exiter interface
func (app *App) Exit(exit Exiter) {
	app.onExit = exit
}

// ExitFunc init func program global var
func (app *App) ExitFunc(f ExiterFunc) {
	app.onExit = f
}

// Thread interface
func (app *App) Thread(threads ...Threader) {
	for _, thread := range threads {
		app.onThreads = append(app.onThreads, thread)
	}
}

// ThreadFunc Func init func program global var
func (app *App) ThreadFunc(fs ...ThreaderFunc) {
	for _, f := range fs {
		app.onThreads = append(app.onThreads, f)
	}
}

func (app *App) runThreads(ctx context.Context) error {
	for id, thread := range app.onThreads {
		if thread == nil {
			return log.Errorf("thread must not nil!")
		}
		go func(ctx context.Context, id int, execute Threader) {
			defer func() {
				if err := recover(); err != nil {
					log.Fatal("Thread[%v] %v\n%v", id, err, string(debug.Stack()))
				}
			}()
			if err := execute.Thread(ctx); err != nil {
				log.Error("Thread[%v]: %v", id, err)
			}
		}(ctx, id, thread)
	}
	return nil
}

// Loop interface
func (app *App) Loop(loops ...Looper) {
	for _, loop := range loops {
		app.onThreadLoops = append(app.onThreadLoops, loop)
	}
}

// LoopFunc Func init func program global var
func (app *App) LoopFunc(fs ...LooperFunc) {
	for _, f := range fs {
		app.onThreadLoops = append(app.onThreadLoops, f)
	}
}

func (app *App) runThreadLoops(ctx context.Context) error {
	for id, threadLoop := range app.onThreadLoops {
		if threadLoop == nil {
			return log.Errorf("threadLoop must not nil!")
		}
		go func(ctx context.Context, id int, execute Looper) {
			thread := func(ctx context.Context, execute Looper) (bool, error) {
				defer func() {
					if err := recover(); err != nil {
						log.Fatal("ThreadLoop[%v] %v\n%v", id, err, string(debug.Stack()))
					}
				}()
				return execute.Loop(ctx)
			}
			var exit bool
			var err error
			for !exit {
				select {
				case <-ctx.Done():
					log.Error("ThreadLoop[%v]: %v", id, ctx.Err())
					return
				default:
					if exit, err = thread(ctx, execute); err != nil {
						log.Error("ThreadLoop[%v]: %v", id, err)
					}
				}
			}
		}(ctx, id, threadLoop)
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

	ok, err := PathExist(app.Option.ConfigPath)
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
	}
	app.Option.Merge(envOpt, argsOpt) // 修改被配置文件改掉配置
	if app.Name == "" {
		app.Name = filepath.Base(os.Args[0])
	}
	TimeZone, err = time.LoadLocation(app.Option.Tz)
	return err
}

// Run app instance
func (app *App) Run(ctx context.Context) {
	if err := app.execute(ctx); err != nil {
		log.Error("%v", err)
	}
}

func (app *App) execute(ctx context.Context) error {
	var err error
	now := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.init()
	defer func(now time.Time) {
		if err := recover(); err != nil {
			log.Fatal("%v\n%v", err, string(debug.Stack()))
		}
		log.Warn("exit: running time=%v", time.Since(now))
	}(now)

	if app.onPrepare != nil {
		if err = app.onPrepare.Prepare(ctx); err != nil {
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

	if err := app.initLog(); err != nil {
		return err
	}
	log.Info("[%s] run...", app.Name)

	if app.onInit != nil {
		if err := app.onInit.Init(ctx); err != nil {
			return err
		}
	}
	log.Warn("%v", string(log.Dump(app)))

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

		mem, err := os.Create(filepath.Join(app.Option.PprofPath, "pprof/mem.prof"))
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
	mainExit := make(chan struct{})
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("MainThread: %v\n%v", err, string(debug.Stack()))
			}
			mainExit <- struct{}{}
		}()
		if app.onMain != nil {
			if err := app.onMain.Main(ctx); err != nil {
				log.Error("MainThread: %v", err)
			}
		}
	}(ctx)

	if err := app.runThreads(ctx); err != nil {
		return err
	}
	if err := app.runThreadLoops(ctx); err != nil {
		return err
	}

	signSystem := make(chan os.Signal)
	signal.Notify(signSystem, os.Interrupt, os.Kill)

	select {
	case <-mainExit:
	case s := <-signSystem:
		log.Warn("Get signal: %v", s)
	case <-ctx.Done():
		log.Warn("%v", ctx.Err())
	}
	if app.onExit != nil {
		return app.onExit.Exit(ctx)
	}
	return nil
}
