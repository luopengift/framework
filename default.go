package framework

import (
	"context"
	"net/http"

	"github.com/luopengift/gohttp"
)

var app *App

func init() {
	app = New()
	app.Regist(NewLog())
}

// Instance return app instance
func Instance() *App {
	return app
}

// WithContext with context
func WithContext(ctx context.Context) {
	app.Context = ctx
}

// BindConfig bind config
func BindConfig(v ConfigProvider) {
	app.BindConfig(v)
}

// InitLimitGroup InitLimitGroup
func InitLimitGroup(max int) {
	app.InitLimitGroup(max)
}

// AppendLimitFunc AppendLimitFunc
func AppendLimitFunc(f FuncVars, vars ...interface{}) {
	app.AppendLimitFunc(f, vars...)
}

// SetPrepareFunc prepare func
func SetPrepareFunc(f PrepareFunc) {
	app.SetPrepareFunc(f)
}

// SetInitFunc init func
func SetInitFunc(f InitFunc) {
	app.SetInitFunc(f)
}

// SetMainFunc main func
func SetMainFunc(f MainFunc) {
	app.SetMainFunc(f)
}

// SetThreadFunc thread func
func SetThreadFunc(f ThreadFunc) {
	app.SetThreadFunc(f)
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

// mainThread default main thread
func mainThread(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func defaultFunc(ctx context.Context) error {
	return nil
}

// SetLogger SetLogger
func SetLogger(provider LogProvider) {
	app.SetLogger(provider)
}

// SetLogProvider SetLogProvider
func SetLogProvider(provider LogProvider) {
	app.SetLogProvider(provider)
}

// Debugf Debugf
func Debugf(format string, v ...interface{}) {
	app.Log.Debugf(format, v...)
}

// Infof Infof
func Infof(format string, v ...interface{}) {
	app.Log.Infof(format, v...)
}

// Warnf warnf
func Warnf(format string, v ...interface{}) {
	app.Log.Warnf(format, v...)
}

// Errorf Errorf
func Errorf(format string, v ...interface{}) {
	app.Log.Errorf(format, v...)
}

// Fatalf Fatalf
func Fatalf(format string, v ...interface{}) {
	app.Log.Fatalf(format, v...)
}

// Regist regist
func Regist(regist interface{}, configs ...interface{}) {
	app.Regist(regist, configs...)
}
