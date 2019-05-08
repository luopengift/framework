package main

import (
	"context"
	"fmt"

	"github.com/luopengift/framework"
	"github.com/luopengift/requests"
)

type Log struct {
	Path string
}
type Config struct {
	Path string `yaml:"path"`
	Log  *framework.Logger
}

func main() {
	sess := requests.New().SetRetry(2).SetTimeout(3)
	c := &Config{
		Path: "ooo",
		Log: &framework.Logger{
			Path: "xxx",
		},
	}
	req := requests.NewRequest("GET", "http://httpbin.org", nil)
	framework.BindConfig(c)
	fmt.Println("---logger", c.Log)
	framework.Regist(framework.NewLog(), c.Log)
	framework.SetMainFunc(func(ctx context.Context) error {
		sess.LogFunc = framework.Instance().Log.Warnf
		if _, err := sess.DoRequest(req); err != nil {
			return err
		}
		// fmt.Println("configProvider", framework.Instance().ConfigProvider.(*Config).Log)
		framework.Infof("%#v", framework.Instance().ConfigProvider)
		framework.Infof("%#v", framework.Instance().Log.LogProvider)

		return nil
	})

	framework.Run()
}
