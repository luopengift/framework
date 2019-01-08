package framework

import (
	"net/http"

	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpd = gohttp.Init()

func (app *App) initHttpd() {
	if app.Option.Httpd == emptyOption.Httpd {
		return
	}
	httpd.Log = log.GetLogger("__ROOT__")
	httpd.RouteStdHandler("^/metrics$", promhttp.Handler())
	go httpd.Run(app.Option.Httpd)
}

// RouteStdHandler route http
func (app *App) RouteStdHandler(path string, handler http.Handler) {
	httpd.RouteStdHandler(path, handler)
}

// Route route http
func (app *App) Route(path string, handler gohttp.Handler) {
	httpd.Route(path, handler)
}

// RouteFunc route http
func (app *App) RouteFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	httpd.RouteFunc(path, f)
}
