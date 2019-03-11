package requests

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request request
type Request struct {
	*http.Request
}

// NewRequest new request
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{req}, nil
}

// WarpRequest warp request
func WarpRequest(req *http.Request) *Request {
	return &Request{req}
}

// StdLib return net/http.Request
func (req *Request) StdLib() *http.Request {
	return req.Request
}

// Query add query args
func (req *Request) Query(query map[string]interface{}) *Request {
	var raw []string
	for k, v := range query {
		raw = append(raw, k+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
	}
	req.Request.URL.RawQuery = strings.Join(raw, "&")
	return req
}

// WithContext returns a shallow copy of r with its context changed to ctx
func (req *Request) WithContext(ctx context.Context) *Request {
	return &Request{req.Request.WithContext(ctx)}
}

// Body request body
func (req *Request) Body(body interface{}) *Request {
	//req.Request.Body =
	return req
}

// SetHeader header
func (req *Request) SetHeader(k, v string) *Request {
	if k != "" && v != "" {
		req.Request.Header.Set(k, v)
	}
	return req
}

// SetHeaders headers
func (req *Request) SetHeaders(kv map[string]string) *Request {
	for k, v := range kv {
		req.SetHeader(k, v)
	}
	return req
}

// Cookie cookie
func (req *Request) Cookie(k, v string) *Request {
	req.Request.AddCookie(&http.Cookie{Name: k, Value: v})
	return req
}

// BaseAuth base auth
func (req *Request) BaseAuth(user, pass string) *Request {
	req.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))
	return req
}
