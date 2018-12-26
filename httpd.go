package framework

import (
	"net/http"

	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpd *gohttp.Application

func (app *App) initHttpd() {
	if app.Option.Httpd == emptyOption.Httpd {
		return
	}
	httpd = gohttp.Init()
	httpd.Log = log.GetLogger("__ROOT__")
	httpd.RouteStdHandler("^/metrics$", promhttp.Handler())
	go httpd.Run(app.Option.Httpd)
}

// RouteStdHandler route http
func (app *App) RouteStdHandler(path string, handler http.Handler) {
	if httpd == nil {
		panic("httpd must init")
	}
	httpd.RouteStdHandler(path, handler)
}

// Route route http
func (app *App) Route(path string, handler gohttp.Handler) {
	if httpd == nil {
		panic("httpd must init")
	}
	httpd.Route(path, handler)
}

// RouteFunc route http
func (app *App) RouteFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	if httpd == nil {
		panic("httpd must init")
	}
	httpd.RouteFunc(path, f)
}
