package main

import (
	"context"

	"github.com/luopengift/framework"
	"github.com/luopengift/requests"
)

type Log struct {
	Path string
}
type Config struct {
	Path string `yaml:"path"`
	Log  Log
}

func main() {
	sess := requests.New().SetRetry(2).SetTimeout(3)
	c := &Config{Path: "ooo"}
	req := requests.NewRequest("GET", "http://httpbin.org", nil)
	app := framework.New()
	app.BindConfig(c)
	// app.Regist(framework.NewLogOpt())
	app.SetMainFunc(func(ctx context.Context) error {
		sess.LogFunc = framework.Instance().Log.Warnf
		if _, err := sess.DoRequest(req); err != nil {
			return err
		}
		app.Log.Infof("%#v", app.Config)
		return nil
	})

	app.Run()
}
