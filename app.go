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
	Config        interface{}
	onPrepare     Prepare
	onInit        Init
	onMain        Main
	onThreads     []Thread
	onThreadLoops []Loop
	onExit        Exiter
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
func (app *App) Parpare(prepare Prepare) {
	app.onPrepare = prepare
}

// ParpareFunc ParpareFunc
func (app *App) ParpareFunc(f PrepareFunc) {
	app.onPrepare = f
}

// Init Initer interface
func (app *App) Init(init Init) {
	app.onInit = init
}

// InitFunc init func program global var
func (app *App) InitFunc(f InitFunc) {
	app.onInit = f
}

// Main interface
func (app *App) Main(main Main) {
	app.onMain = main
}

// MainFunc handle main func
func (app *App) MainFunc(f MainFunc) {
	app.onMain = f
}

// Exit Exiter interface
func (app *App) Exit(exit Exiter) {
	app.onExit = exit
}

// ExitFunc init func program global var
func (app *App) ExitFunc(f ExitFunc) {
	app.onExit = f
}

// Thread interface
func (app *App) Thread(threads ...Thread) {
	for _, thread := range threads {
		app.onThreads = append(app.onThreads, thread)
	}
}

// ThreadFunc Func init func program global var
func (app *App) ThreadFunc(fs ...ThreadFunc) {
	for _, f := range fs {
		app.onThreads = append(app.onThreads, f)
	}
}

func (app *App) runThreads(ctx context.Context) error {
	for idx, thread := range app.onThreads {
		if thread == nil {
			return log.Errorf("thread must not nil!")
		}
		go func(id int, execute Thread, ctx context.Context, app *App) {
			if err := execute.Thread(ctx, app); err != nil {
				log.Error("Thread[%v]: %v", id, err)
			}
		}(idx, thread, ctx, app)
	}
	return nil
}

// Loop interface
func (app *App) Loop(loops ...Loop) {
	for _, loop := range loops {
		app.onThreadLoops = append(app.onThreadLoops, loop)
	}
}

// LoopFunc Func init func program global var
func (app *App) LoopFunc(fs ...LoopFunc) {
	for _, f := range fs {
		app.onThreadLoops = append(app.onThreadLoops, f)
	}
}

func (app *App) runThreadLoops(ctx context.Context) error {
	for idx, threadLoop := range app.onThreadLoops {
		if threadLoop == nil {
			return log.Errorf("threadLoop must not nil!")
		}
		go func(id int, execute Loop, ctx context.Context, app *App) {
			var exit bool
			var err error
			for !exit {
				select {
				case <-ctx.Done():
					log.Error("ThreadLoop[%v]: %v", id, ctx.Err())
					return
				default:
					if exit, err = execute.Loop(app); err != nil {
						log.Error("ThreadLoop[%v]: %v", id, err)
					}
				}
			}
		}(idx, threadLoop, ctx, app)
	}
	return nil
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

	if app.onMain == nil {
		return log.Errorf("Main is nil, must set!")
	}
	signExit := make(chan struct{})
	go func(ctx context.Context, app *App) {
		if err := app.onMain.Main(ctx, app); err != nil {
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
