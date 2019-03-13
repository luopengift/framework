package main

import (
	"github.com/luopengift/framework"
)

type config struct {
	Httpd string
}

func main() {
	framework.BindConfig(&config{
		Httpd: ":9099",
	})
	framework.HttpdRoute("/report", &framework.Report{})
	framework.Run()
}
