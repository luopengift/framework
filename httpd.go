package framework

import (
	"github.com/luopengift/gohttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func httpd(addr string) error {
	app := gohttp.Init()
	app.RouteStdHandler("/metrics", promhttp.Handler())
	app.Run(addr)
	return nil
}
