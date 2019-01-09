package framework

import (
	"context"
	"net/http"

	"github.com/luopengift/gohttp"
)

var app *App

func init() {
	app = New()
	app.MainFunc(defaultMainThread)
}

// WithContext with context
func WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// Bind bind runner interface
func Bind(r Runner) {
	app.Bind(r)
}

// BindConfig bind config
func BindConfig(v interface{}) {
	app.BindConfig(v)
}

// Prepare prepare interface
func Prepare(v Preparer) {
	app.Prepare(v)
}

// PrepareFunc prepare func
func PrepareFunc(f PreparerFunc) {
	app.PrepareFunc(f)
}

// Init init interface
func Init(v Initer) {
	app.Init(v)
}

// InitFunc init func
func InitFunc(f IniterFunc) {
	app.InitFunc(f)
}

// Main main interface
func Main(v Mainer) {
	app.Main(v)
}

// MainFunc main func
func MainFunc(f MainerFunc) {
	app.MainFunc(f)
}

// Thread thread interface
func Thread(v ...Threader) {
	app.Thread(v...)
}

// ThreadFunc thread func
func ThreadFunc(f ...ThreaderFunc) {
	app.ThreadFunc(f...)
}

// Loop loop interface
// func Loop(v ...Looper) {
// 	app.Loop(v...)
// }

// LoopFunc loop func
// func LoopFunc(f ...LooperFunc) {
// 	app.LoopFunc(f...)
// }

// GoroutineFunc GoroutineFunc
func GoroutineFunc(name string, v GoroutinerFunc, num ...int) {
	app.GoroutineFunc(name, v, num...)
}

// Exit interface
func Exit(v Exiter) {
	app.Exit(v)
}

// ExitFunc exit func
func ExitFunc(f ExiterFunc) {
	app.ExitFunc(f)
}

// HttpdRoute http route
func HttpdRoute(path string, handler gohttp.Handler) {
	app.Route(path, handler)
}

// HttpdRouteFunCtx http rout func ctx
func HttpdRouteFunCtx(path string, handler gohttp.HandleFunCtx) {
	app.RouteFunCtx(path, handler)
}

// HttpdRouteFunc http route func
func HttpdRouteFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	app.RouteFunc(path, f)
}

// Run run
func Run(ctx ...context.Context) {
	var c context.Context
	if len(ctx) == 0 {
		c = context.Background()
	}
	app.Run(c)
}

// NewReport new report
func NewReport() *Report {
	return app.NewReport()
}

// defaultMainThread default main thread
func defaultMainThread(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
