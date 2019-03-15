package framework

import (
	"context"
	"net/http"

	"github.com/luopengift/gohttp"
)

var app *App

func init() {
	app = New()
}

// Instance return app instance
func Instance() *App {
	return app
}

// WithContext with context
func WithContext(ctx context.Context) {
	app.Context = ctx
}

// Bind bind runner interface
func Bind(r Runner) {
	app.Bind(r)
}

// BindConfig bind config
func BindConfig(v interface{}) {
	app.BindConfig(v)
}

// InitLog  init Log module
func InitLog() error {
	return app.InitLog()
}

// LoadConfig load config
func LoadConfig() error {
	return app.LoadConfig()
}

// InitLimitGroup InitLimitGroup
func InitLimitGroup(max int) {
	app.InitLimitGroup(max)
}

// AppendLimitFunc AppendLimitFunc
func AppendLimitFunc(f FuncVars, vars ...interface{}) {
	app.AppendLimitFunc(f, vars...)
}

// Prepare prepare interface
func Prepare(v Function) {
	app.Prepare(v)
}

// PrepareFunc prepare func
func PrepareFunc(f Func) {
	app.PrepareFunc(f)
}

// Init init interface
func Init(v Function) {
	app.Init(v)
}

// InitFunc init func
func InitFunc(f Func) {
	app.InitFunc(f)
}

// Main main interface
func Main(v Function) {
	app.Main(v)
}

// MainFunc main func
func MainFunc(f Func) {
	app.MainFunc(f)
}

// Thread thread interface
func Thread(v ...Function) {
	app.Thread(v...)
}

// ThreadFunc thread func
func ThreadFunc(f ...Func) {
	app.ThreadFunc(f...)
}

// GoroutineFunc GoroutineFunc
func GoroutineFunc(name string, v FuncWithExit, num ...int) {
	app.GoroutineFunc(name, v, num...)
}

// Exit interface
func Exit(v Function) {
	app.Exit(v)
}

// ExitFunc exit func
func ExitFunc(f Func) {
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
func Run() {
	app.Run()
}

// NewReport new report
func NewReport() *Report {
	return app.NewReport()
}

// DefaultMainThread default main thread
func DefaultMainThread(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func defaultFunc(ctx context.Context) error {
	return nil
}
