package framework

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *App) initHttpd() {
	if app.Option.Httpd != emptyOption.Httpd {
		app.Application = gohttp.Init()
		app.Application.Log = log.GetLogger("__ROOT__")
		app.Application.RouteStdHandler("/metrics", promhttp.Handler())
		go app.Application.Run(app.Option.Httpd)
	}
}
