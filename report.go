package framework

// Report report
type Report struct {
	URLs []string
}

func (r *Report) getURL() string {
	return ""
}

func (r *Report) request() {
	//http.Post(r.URLs[0])
}
