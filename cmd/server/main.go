package main

import (
	"github.com/luopengift/framework"
)

func main() {
	framework.HttpdRoute("/report", &framework.Report{})
	framework.Run()
}
