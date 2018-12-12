package framework

import (
	"context"
	"flag"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

type Executer interface {
	Execute(ctx context.Context, app *App) error
}

type Func func(ctx context.Context, app *App) error

func (f Func) Execute(ctx context.Context, app *App) error {
	return f(ctx, app)
}

// App framework
type App struct {
	*Option
	Config interface{}
	Func   Executer
}

// Option option
type Option struct {
	Debug   bool
	LogPath string
}

func (app *App) copy(options ...*App) {
	for _, option := range options {
		if option.Debug != false {
			app.Debug = option.Debug
		}
		if option.Config != nil {
			app.Config = option.Config
		}
		if option.Func != nil {
			app.Func = option.Func
		}
	}
}

// New new app instance
func New() *App {
	return &App{
		Option: &Option{},
	}

}

func (app *App) HandleFunc(f Func) {
	app.Func = f
}

// Init app instance
func (app *App) Init() error {
	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	d := flag.Bool("debug", false, "(debug)调试模式")
	v := flag.Bool("version", false, "(version)版本")
	//addr := flag.String("http", ":8888", "(http)地址")
	flag.Parse()
	if *v {
		log.ConsoleWithMagenta("%v", version.String())
		return io.EOF
	}
	app.Debug = *d
	return types.ParseConfigFile(app.Config, *c)
}

// Run app instance
func (app *App) Run(ctx context.Context) {
	now := time.Now()
	done := make(chan struct{})
	log.Display("CONFIG", app)
	file := log.NewFile("%Y-%M-%D.log")
	file.SetMaxBytes(1000 * 1024 * 1024)
	if app.Debug {
		log.SetLevel(log.DEBUG)
		log.SetOutput(file, os.Stderr)
	} else {
		log.SetLevel(log.INFO)
		log.SetOutput(file)
	}

	if app.Func == nil {
		log.Error("Func must set! %#v", &app.Func)
		return
	}
	go func(ctx context.Context, app *App) {
		if err := app.Func.Execute(ctx, app); err != nil {
			log.Error("Execute: %v", err)
		}
		done <- struct{}{}
	}(ctx, app)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	select {
	case <-sign:
		log.Warn("[CTRL+C]")
	case <-done:
	case <-ctx.Done():
	}
	log.Info("running time: %v", time.Since(now))
}
