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

var mainFunc = func(ctx context.Context, app *framework.App) error {
	time.Sleep(10 * time.Second)
	return nil
}
var t1 = func(app *framework.App) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	return false, nil
}

func main() {
	ctx := context.Background()
	app := framework.New()
	//app.BindConfig(&config{})

	app.MainFunc(mainFunc)
	app.LoopFunc(t1)
	if err := app.Run(ctx); err != nil {
		log.Error("%v", err)
	}
}
