package main

import (
	"context"

	"github.com/luopengift/framework"
	"github.com/luopengift/log"
)

var F = func(ctx context.Context, app *framework.App) error {
	log.Debug("debug...")
	return nil
}

func main() {
	ctx := context.Background()
	app := framework.New()
	app.Init()
	app.Debug = true

	app.HandleFunc(F)
	app.Run(ctx)
}
