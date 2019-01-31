package main

import (
	"context"

	"github.com/luopengift/framework"
)

type config struct {
	ReportURL string `json:"report_url" yaml:"report_url"`
}

type run struct{}

// Prepare prepare
func (r *run) Prepare(ctx context.Context) error {
	return nil
}

// Init init
func (r *run) Init(ctx context.Context) error {
	return nil
}

func (r *run) Main(ctx context.Context) error {
	//	time.Sleep(60 * time.Second)
	//select {}
	return nil
}

func (r *run) Thread(ctx context.Context) error {
	//panic("ddd")
	return nil
}

func (r *run) Exit(ctx context.Context) error {
	return nil
}

func main() {
	//framework.
	framework.Bind(&run{})
	framework.GoroutineFunc("", func(ctx context.Context) (bool, error) {
		return false, nil
	}, 10)
	framework.Run()
}
