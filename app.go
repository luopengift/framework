package framework

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/luopengift/framework/pkg/limit"
	"github.com/luopengift/framework/pkg/log"
	"github.com/luopengift/framework/util"
)

// var default var
var (
	sTime   time.Time
	wg      sync.WaitGroup
	mux     sync.RWMutex
	one     sync.Once
	limiter *limit.Limit
)

// App framework
type App struct {
	context.Context `json:"-"`
	Name            string `json:"name" yaml:"name"`
	ID              string `json:"id" yaml:"id"`
	Log             *Log
	*Option         `json:"option"`
	TimeZone        *time.Location `json:"timeZone"`
	Config          ConfigProvider
	*run
	raws     []byte
	register []interface{}
}

// NewOnce new app by once
func NewOnce(opts ...interface{}) *App {
	mux.Lock()
	defer mux.Unlock()
	one.Do(func() {
		app = New(opts...)
	})
	return app
}

// New new app instance
func New(opts ...interface{}) *App {
	sTime = time.Now()
	app := &App{
		Context: context.Background(),
		ID:      util.Random(10),
		Name:    filepath.Base(os.Args[0]),
		Log: &Log{
			log.NewStdLog(),
		},
		Config: struct{}{},
		Option: defaultOption,
		run:    &run{},
	}
	app.SetPrepareFunc(defaultFunc)
	app.SetInitFunc(defaultFunc)
	app.SetMainFunc(mainThread)
	app.SetExitFunc(defaultFunc)
	app.SetOpts(opts...)
	return app
}

// SetOpts set opts
func (app *App) SetOpts(opts ...interface{}) *App {
	for _, opt := range opts {
		switch v := opt.(type) {
		case LogProvider:
			app.Log.LogProvider = v
		case *Option:
			app.mergeIn(v)
		case Option:
			app.mergeIn(&v)
		default:
			panic("Unknow type opts")
		}
	}
	return app
}

// WithContext with context
func (app *App) WithContext(ctx context.Context) {
	app.Context = ctx
}

// BindConfig bind config
func (app *App) BindConfig(v ConfigProvider) {
	app.Config = v
}

// InitLimitGroup initLimitGroup
func (app *App) InitLimitGroup(max int) {
	limiter = limit.NewLimit(max)
}

// FuncVars func with vars
type FuncVars func(ctx context.Context, vars ...interface{}) error

// AppendLimitFunc  append limitFunc
func (app *App) AppendLimitFunc(f FuncVars, vars ...interface{}) {
	limiter.Add()
	go func() {
		id := util.Random(10)
		defer func() {
			limiter.Done()
			if err := recover(); err != nil {
				app.PrintStack(app.Log.Fatalf, "Limit panic[%v]: %v\n%v", id, err)
			}
		}()
		if err := f(app.Context, vars...); err != nil {
			app.PrintStack(app.Log.Errorf, "Limit[%v]: %v", id, err)
		}
	}()
}

// PrintStack append stack with func(format string ,v ...interface{})
func (app *App) PrintStack(f func(string, ...interface{}), format string, v ...interface{}) {
	f(format+"\n%s", append(v, string(debug.Stack()))...)
}

// WaitLimitDone Wait
func (app *App) WaitLimitDone() {
	limiter.Wait()
}

// Run app instance
func (app *App) Run() {
	defer func() {
		if err := recover(); err != nil {
			app.PrintStack(app.Log.Errorf, "%v", err)
		}
		app.Log.Warnf("exit: running time=%v", time.Since(sTime))
	}()
	if err := app.execute(); err != nil {
		app.PrintStack(app.Log.Errorf, "%v", err)
	}
	if err := app.onExit.ExitFunc(app.Context); err != nil {
		app.PrintStack(app.Log.Errorf, "onExit: %v", err)
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

	if err = app.onPrepare.PrepareFunc(ctx); err != nil {
		return err
	}

	if err = app.LoadConfig(); err != nil {
		return err
	}

	app.SetOpts(app.register...)

	if err = app.onInit.InitFunc(ctx); err != nil {
		return err
	}
	app.Log.Infof("[%s] init...", app.Name)

	// http
	app.initHttpd()

	// pprof
	if app.Option.PprofPath != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = startPprof(ctx, app.Option.PprofPath); err != nil {
				app.Log.Errorf("%v", err)
			}
		}()
	}

	app.Log.Infof("init success. cost=%v", time.Since(sTime))

	mainExit := make(chan error)
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				app.PrintStack(app.Log.Fatalf, "MainThread: %v", err)
				mainExit <- fmt.Errorf("%v", err)
			}
		}()
		mainExit <- app.onMain.MainFunc(ctx)
	}(ctx)

	if err = app.runThreads(ctx); err != nil {
		return err
	}

	if err = app.runGoroutines(ctx); err != nil {
		return err
	}

	signSystem := make(chan os.Signal)
	signal.Notify(signSystem, os.Interrupt, os.Kill)

	select {
	case err = <-mainExit:
		app.Log.Warnf("Exit MainThread: %v", err)
	case s := <-signSystem:
		app.Log.Warnf("Exit: signal %v", s)
	}
	return nil
}

// Regist v into framework
func (app *App) Regist(v interface{}) {
	// return json.Unmarshal(app.raws, v)
	app.register = append(app.register, v)

}

// func (app *App) RegistLog(provider LogProvider) error {
// 	app.Log.LogProvider = provider
// }
