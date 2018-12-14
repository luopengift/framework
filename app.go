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

// App framework
type App struct {
	*Option
	exit        chan struct{} //退出信号
	Config      interface{}
	onInit      func(*App) error
	onFlag      func() error
	Func        Executer
	Threads     []Executer
	ThreadLoops []ExecuteLooper
}

// New new app instance
func New(opts ...*Option) *App {
	app := &App{
		Option: defaultOption,
		exit:   make(chan struct{}),
	}
	app.Option.Merge(opts...)
	return app
}

// Flag flag
func (app *App) Flag(fun func() error) {
	app.onFlag = fun
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
// func (app *App) Init(opts ...*Option) error {
// 	app.Option.Merge(opts...)
// 	return nil
// }
// Init init program global var
func (app *App) Init(fun func(app *App) error) {
	app.onInit = fun
}

// Run app instance
func (app *App) Run(ctx context.Context) error {
	if app.onFlag != nil {
		if err := app.onFlag(); err != nil {
			return err
		}
	}
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
	if app.onInit != nil {
		if err := app.onInit(app); err != nil {
			return err
		}
	}
	flag.Parse()
	if app.Func == nil {
		return log.Errorf("Func must set! %T", app.Func.Execute)
	}

	go func(ctx context.Context, app *App) {
		if err := app.Func.Execute(ctx, app); err != nil {
			log.Error("Execute: %v", err)
		}
		app.exit <- struct{}{}
	}(ctx, app)

	if err := app.runThreads(ctx); err != nil {
		return err
	}
	if err := app.runThreadLoops(ctx); err != nil {
		return err
	}

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
