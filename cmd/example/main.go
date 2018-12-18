package main

import (
	"context"
	"fmt"
	"time"

	"github.com/luopengift/framework"
	"github.com/luopengift/log"
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
	time.Sleep(10 * time.Second)
	return nil
}

func (r *run) Thread(ctx context.Context) error {
	return nil
}

func (r *run) Loop(ctx context.Context) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	return false, nil
}

func (r *run) Exit(ctx context.Context) error {
	return nil
}

func main() {
	ctx := context.Background()
	app := framework.New()
	//app.BindConfig(&config{})
	app.Bind(&run{})
	if err := app.Run(ctx); err != nil {
		log.Error("%v", err)
	}
}
