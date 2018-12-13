package framework

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

// Executer execute interface, 用户自行管理携程运行退出等状态, framework仅调起函数
type Executer interface {
	Execute(ctx context.Context, app *App) error
}

// Func func
type Func func(ctx context.Context, app *App) error

// Execute execute
func (f Func) Execute(ctx context.Context, app *App) error {
	return f(ctx, app)
}

// ExecuteLooper 协程循环由framework管理
type ExecuteLooper interface {
	ExecuteLoop(app *App) (bool, error)
}

// FuncLoop func
type FuncLoop func(app *App) (bool, error)

// ExecuteLoop ExecuteLoop
func (f FuncLoop) ExecuteLoop(app *App) (bool, error) {
	return f(app)
}

// App framework
type App struct {
	*Option
	exit        chan struct{} //退出信号
	Config      interface{}
	Func        Executer
	Threads     []Executer
	ThreadLoops []ExecuteLooper
}

var emptyOption = &Option{}
var defaultOption = &Option{
	Debug:          true,
	LogPath:        "logs/%Y-%M-%D.log",
	MaxBytes:       200 * 1024 * 1024, //200M
	MaxBackupIndex: 50,
}

// Option option
type Option struct {
	Debug          bool
	LogPath        string
	MaxBytes       int // 日志文件大小
	MaxBackupIndex int // 日志文件数量
}

func (opt *Option) mrgreIn(o *Option) {
	if o.Debug != emptyOption.Debug {
		opt.Debug = o.Debug
	}
	if o.LogPath != emptyOption.LogPath {
		opt.LogPath = o.LogPath
	}
	if o.MaxBytes != emptyOption.MaxBytes {
		opt.MaxBytes = o.MaxBytes
	}
	if o.MaxBackupIndex != emptyOption.MaxBackupIndex {
		opt.MaxBackupIndex = o.MaxBackupIndex
	}
}

// Merge merge
func (opt *Option) Merge(opts ...*Option) {
	for _, o := range opts {
		opt.mrgreIn(o)
	}
}

// New new app instance
func New() *App {
	return &App{
		Option: defaultOption,
		exit:   make(chan struct{}),
	}
}

// HandleFunc handle func
func (app *App) HandleFunc(f Func) {
	app.Func = f
}

// MainLoopFunc main loop func
func (app *App) MainLoopFunc(f Func) {
	app.Func = f
}

// ThreadFuncs thread funcs
func (app *App) ThreadFuncs(funs ...Func) {
	for _, fun := range funs {
		app.Threads = append(app.Threads, fun)
	}
}

func (app *App) runThreads(ctx context.Context) error {
	for idx, thread := range app.Threads {
		if thread == nil {
			return log.Errorf("thread must not nil!")
		}
		go func(id int, execute Executer, ctx context.Context, app *App) {
			if err := execute.Execute(ctx, app); err != nil {
				log.Error("Thread[%v]: %v", id, err)
			}
		}(idx, thread, ctx, app)
	}
	return nil
}

func (app *App) runThreadLoops(ctx context.Context) error {
	for idx, threadLoop := range app.ThreadLoops {
		if threadLoop == nil {
			return log.Errorf("threadLoop must not nil!")
		}
		go func(id int, execute ExecuteLooper, ctx context.Context, app *App) {
			var exit bool
			var err error
			for !exit {
				select {
				case <-ctx.Done():
					log.Error("ThreadLoop[%v]: %v", id, ctx.Err())
					return
				default:
					if exit, err = execute.ExecuteLoop(app); err != nil {
						log.Error("ThreadLoop[%v]: %v", id, err)
					}
				}
			}
		}(idx, threadLoop, ctx, app)
	}
	return nil
}

// ThreadLoopFuncs thread loop funcs
func (app *App) ThreadLoopFuncs(funs ...FuncLoop) {
	for _, fun := range funs {
		app.ThreadLoops = append(app.ThreadLoops, fun)
	}
}

// Init app instance
func (app *App) Init(opts ...*Option) error {
	app.Option.Merge(opts...)
	return nil
}

// Run app instance
func (app *App) Run(ctx context.Context) error {
	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	//p := flag.Bool("pprof", false, "(pprof)调试模式")
	v := flag.Bool("version", false, "(version)版本")
	//addr := flag.String("http", ":8888", "(http)地址")
	flag.Parse()
	if *v {
		log.ConsoleWithMagenta("%v", version.String())
		return nil
	}

	now := time.Now()
	defer func(now time.Time) {
		log.Warn("[EXIT]running time=%v", time.Since(now))
	}(now)

	if app.Config != nil {
		if err := types.ParseConfigFile(app.Config, *c); err != nil {
			return err
		}
	}
	if err := app.initLog(); err != nil {
		return err
	}
	if app.Func == nil {
		return log.Errorf("Func must set! %T", app.Func.Execute)
	}

	if err := app.runThreads(ctx); err != nil {
		return err
	}
	if err := app.runThreadLoops(ctx); err != nil {
		return err
	}

	go func(ctx context.Context, app *App) {
		if err := app.Func.Execute(ctx, app); err != nil {
			log.Error("Execute: %v", err)
		}
		app.exit <- struct{}{}
	}(ctx, app)
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill)
	select {
	case <-app.exit:
	case <-sign:
		log.Warn("[CTRL+C]")
	case <-ctx.Done():
		log.Warn("%v", ctx.Err())
	}
	return nil
}
