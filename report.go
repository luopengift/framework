package framework

import "time"

// Report report
type Report struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

func (r *Report) getURL() string {
	return ""
}

func (r *Report) request() {
	//http.Post(r.URLs[0])
}
