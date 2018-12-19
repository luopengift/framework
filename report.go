package framework

import (
	"encoding/json"
	"time"

	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
)

// Report report
type Report struct {
	gohttp.BaseHTTPHandler `json:"-"`
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Time                   time.Time `json:"time"`
}

// NewReport new report instance
func NewReport(name string) *Report {
	return &Report{
		ID:   "xxxx-xxxx-xxxx",
		Name: name,
		Time: time.Now(),
	}
}

// GET method
func (r *Report) GET() {
	log.Info("/report get...")
}

// POST method
func (r *Report) POST() {
	if err := json.Unmarshal(r.GetBodyArgs(), r); err != nil {
		log.Error("%v", err)
		return
	}
	log.Info("%v", string(r.GetBodyArgs()))
	r.Output(r)
}

func (r *Report) getURL() string {
	return ""
}

func (r *Report) request() {
	//http.Post(r.URLs[0])
}
