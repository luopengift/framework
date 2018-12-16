package framework

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

// App framework
type App struct {
	*Option
	Config      interface{}
	onPrepare   Preparer
	onInit      Initer
	Main        Executer
	Threads     []Executer
	ThreadLoops []ExecuteLooper
	onExit      Exiter
}

// New new app instance
func New(opts ...*Option) *App {
	app := &App{
		Option: defaultOption,
	}
	app.Option.Merge(opts...)
	return app
}

// BindConfig bind config
func (app *App) BindConfig(v interface{}) {
	app.Config = v
}

// Parpare Preparer interface
func (app *App) Parpare(prepare Preparer) {
	app.onPrepare = prepare
}

// ParpareFunc ParpareFunc
func (app *App) ParpareFunc(f PrepareFunc) {
	app.onPrepare = f
}

// Init Initer interface
func (app *App) Init(init Initer) {
	app.onInit = init
}

// InitFunc init func program global var
func (app *App) InitFunc(f InitFunc) {
	app.onInit = f
}

// Exit Exiter interface
func (app *App) Exit(exit Exiter) {
	app.onExit = exit
}

// ExitFunc init func program global var
func (app *App) ExitFunc(f ExitFunc) {
	app.onExit = f
}

// HandleFunc handle func
func (app *App) HandleFunc(f Func) {
	app.Main = f
}

// MainLoopFunc main loop func
func (app *App) MainLoopFunc(f Func) {
	app.Main = f
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

// Run app instance
func (app *App) Run(ctx context.Context) error {
	var err error
	now := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer func(now time.Time) {
		log.Warn("exit: running time=%v", time.Since(now))
	}(now)

	if app.onPrepare != nil {
		if err = app.onPrepare.Prepare(); err != nil {
			return err
		}
	}

	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	d := flag.Bool("debug", false, "(debug)调试模式")
	p := flag.Bool("pprof", false, "(pprof)性能分析")
	v := flag.Bool("version", false, "(version)版本")
	//addr := flag.String("http", ":8888", "(http)地址")
	flag.Parse()
	app.Option.Debug = *d
	if *v {
		log.ConsoleWithMagenta("%v", version.String())
		return nil
	}

	if err := app.initLog(); err != nil {
		return err
	}

	// pprof
	if *p {
		os.MkdirAll("var", 0755)
		cpu, err := os.Create("pprof/cpu.prof")
		if err != nil {
			return err
		}
		defer cpu.Close()

		if err = pprof.StartCPUProfile(cpu); err != nil {
			return err
		}
		defer pprof.StopCPUProfile()

		mem, err := os.Create("pprof/mem.prof")
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

	if app.Config != nil {
		if err := types.ParseConfigFile(app.Config, *c); err != nil {
			return err
		}
		if err := types.ParseConfigFile(app.Option, *c); err != nil {
			return err
		}
	}

	if app.onInit != nil {
		if err := app.onInit.Init(app); err != nil {
			return err
		}
	}

	log.Display("%v", app)

	if app.Main == nil {
		return log.Errorf("Main is nil, must set!")
	}
	signExit := make(chan struct{})
	go func(ctx context.Context, app *App) {
		if err := app.Main.Execute(ctx, app); err != nil {
			log.Error("Execute: %v", err)
		}
		signExit <- struct{}{}
	}(ctx, app)

	if err := app.runThreads(ctx); err != nil {
		return err
	}
	if err := app.runThreadLoops(ctx); err != nil {
		return err
	}

	signSystem := make(chan os.Signal)
	signal.Notify(signSystem, os.Interrupt, os.Kill)

	select {
	case <-signExit:
	case s := <-signSystem:
		log.Warn("Get signal: %v", s)
	case <-ctx.Done():
		log.Warn("%v", ctx.Err())
	}
	if app.onExit != nil {
		return app.onExit.Exit(ctx, app)
	}
	return nil
}
