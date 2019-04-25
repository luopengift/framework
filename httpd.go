package framework

import (
	"net/http"

	"github.com/luopengift/framework/pkg/path/pathutil"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpd = gohttp.Init()

func (app *App) initHttpd() error {
	if app.Option.Httpd == emptyOption.Httpd {
		return nil
	}
	if ok, err := pathutil.Exist(app.Option.ConfigPath); ok && err == nil {
		if err := types.ParseConfigFile(httpd.Config, app.Option.ConfigPath); err != nil {
			return err
		}
	}
	httpd.Log = app.Log
	httpd.RouteStdHandler("^/metrics$", promhttp.Handler())
	go httpd.Run(app.Option.Httpd)
	return nil
}

// RouteStdHandler route http
func (app *App) RouteStdHandler(path string, handler http.Handler) {
	httpd.RouteStdHandler(path, handler)
}

// Route route http
func (app *App) Route(path string, handler gohttp.Handler) {
	httpd.Route(path, handler)
}

// RouteFunCtx route fun ctx
func (app *App) RouteFunCtx(path string, handler gohttp.HandleFunCtx) {
	httpd.RouteFunCtx(path, handler)
}

// RouteFunc route http
func (app *App) RouteFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	httpd.RouteFunc(path, f)
}
