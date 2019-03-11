package requests

import (
	"io"
	"net/url"
)

var session = New()

// Do request
func Do(method, url, contentType string, body io.Reader) (*Response, error) {
	return session.Do(method, url, contentType, body)
}

// Get send get request
func Get(url string) (resp *Response, err error) {
	return session.Get(url)
}

// Post send post request
func Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	return session.Post(url, contentType, body)
}

// PostForm send post request,  content-type = application/x-www-form-urlencoded
func PostForm(url string, data url.Values) (resp *Response, err error) {
	return session.PostForm(url, data)
}
