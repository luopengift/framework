package main

import (
	"context"
	"flag"

	"github.com/luopengift/framework"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
)

type runner struct {
	*framework.Run
	addr  *string
	httpd *gohttp.Application
}

func (r *runner) Prepare(ctx context.Context) error {
	r.addr = flag.String("http", ":3456", "(http)")
	return nil
}

func (r *runner) Init(ctx context.Context) error {
	r.httpd = gohttp.Init()
	r.httpd.Log = log.GetLogger("__ROOT__")
	r.httpd.Route("/report", &framework.Report{})
	return nil
}

func (r *runner) Main(ctx context.Context) error {
	r.httpd.Run(*r.addr)
	return nil
}

func main() {
	ctx := context.Background()
	app := framework.New()
	app.Bind(&runner{})
	if err := app.Run(ctx); err != nil {
		log.Error("%v", err)
	}
}
