package main

import (
	"encoding/json"
	"time"

	"github.com/luopengift/framework"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
)

type config struct {
	Httpd string
}

// Report report
type Report struct {
	gohttp.BaseHTTPHandler `json:"-"`
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Time                   time.Time `json:"time"`
}

// NewReport new report instance
func NewReport() *Report {
	return &Report{
		ID:   "app.ID",
		Name: "app.Name",
		Time: time.Now(),
	}
}

// GET method
func (r *Report) GET() {
	log.Infof("/report get...")
}

// POST method
func (r *Report) POST() {
	if err := json.Unmarshal(r.GetBodyArgs(), r); err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("%v", string(r.GetBodyArgs()))
	r.Output(r)
}

func (r *Report) getURL() string {
	return ""
}

func (r *Report) request() {
	//http.Post(r.URLs[0])
}

func main() {
	framework.BindConfig(&config{
		Httpd: ":9099",
	})
	framework.HttpdRoute("/report", &Report{})
	conf := framework.Instance().Config.(*config)
	log.Infof("%#v", conf.Httpd)
	framework.Run()
}
