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

// App framework
type App struct {
	Debug  bool
	Config interface{}
	Func   func(ctx context.Context, app *App) error
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
	return &App{}
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
	log.Display("CONFIG", app.Config)
	if app.Debug {
		log.SetLevel(log.DEBUG)
	} else {
		log.SetLevel(log.INFO)
	}
	file := log.NewFile("%Y-%M-%D.log")
	file.SetMaxBytes(1000 * 1024 * 1024)
	log.SetOutput(file, os.Stderr)
	go app.Func(ctx, app)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	select {
	case <-sign:
		log.Warn("[CTRL+C]")
	}
	log.Info("running time: %v", time.Since(now))
}
