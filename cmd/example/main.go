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
var t1 = func(ctx context.Context, app *framework.App) error {
	log.Debug("debug...")
	return fmt.Errorf("t1 exit")
}

func main() {
	ctx := context.Background()
	app := framework.New()
	//app.Init(nil)
	//app.Config = &config{}
	app.Debug = true

	app.MainLoopFunc(mainFunc)
	app.ThreadFuncs(t1)
	if err := app.Run(ctx); err != nil {
		log.Error("%v", err)
	}
}
