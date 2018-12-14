package main

import (
	"context"
	"fmt"
	"time"

	"github.com/luopengift/framework"
	"github.com/luopengift/log"
)

type config struct {
	url string
}

var mainFunc = func(ctx context.Context, app *framework.App) error {
	log.Debug("debug...")
	time.Sleep(1 * time.Second)
	return nil
}
var t1 = func(app *framework.App) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	log.Debug("debug...")
	//panic("panic test")
	return false, nil
}

func main() {
	ctx := context.Background()
	app := framework.New()
	//app.Init(nil)
	//app.Config = &config{}
	app.Debug = true
	app.Option.Debug = false

	app.MainLoopFunc(mainFunc)
	app.ThreadLoopFuncs(t1)
	if err := app.Run(ctx); err != nil {
		log.Error("%v", err)
	}
}
