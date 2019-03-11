package requests

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// var
var (
	ErrEmptyProxy = errors.New("proxy is empty")
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"

// Session httpclient session
// Clients and Transports are safe for concurrent use by multiple goroutines
// for efficiency should only be created once and re-used.
// so, session is also safe for concurrent use by multiple goroutines.
type Session struct {
	*http.Client
	*http.Transport
	err chan error
}

// New new session
func New() *Session {
	var (
		errs   = make(chan error)
		tr     = &http.Transport{}
		client = &http.Client{}
	)
	client.Timeout = 120 * time.Second
	client.Transport = tr
	tr.MaxIdleConns = 10
	tr.IdleConnTimeout = 120 * time.Second
	tr.DisableCompression = true
	tr.DisableKeepAlives = false

	return &Session{client, tr, errs}
}

// Proxy set proxy addr
// os.Setenv("HTTP_PROXY", "http://127.0.0.1:9743")
// os.Setenv("HTTPS_PROXY", "https://127.0.0.1:9743")
func (sess *Session) Proxy(addr string) error {
	if addr == "" {
		return ErrEmptyProxy
	}
	proxyURL, err := url.Parse(addr)
	if err != nil {
		return err
	}
	switch proxyURL.Scheme {
	case "socks5", "socks4":
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
		if err != nil {
			return err
		}
		sess.Transport.Dial = dialer.Dial
	default:
		sess.Transport.Proxy = http.ProxyURL(proxyURL)
	}
	return nil
}

// Timeout set client timeout
func (sess *Session) Timeout(timeout int) {
	sess.Client.Timeout = time.Duration(timeout) * time.Second
}

// DisableKeepAlives set transport disableKeepAlives default transport is keepalive,
// if set true, only use the connection to the server for a single HTTP request.
func (sess *Session) DisableKeepAlives(disableKeepAlives bool) {
	sess.Transport.DisableKeepAlives = disableKeepAlives
}

// DoRequest send a request and return a response
func (sess *Session) DoRequest(req *Request) (*Response, error) {
	resp, err := sess.Client.Do(req.Request)
	return &Response{resp}, err
}

// DoRequestWithContext Do with context
func (sess *Session) DoRequestWithContext(ctx context.Context, req *Request) (*Response, error) {
	req2 := req.WithContext(ctx) // !!! WithContext returns a shallow copy of r with its context changed to ctx
	return sess.DoRequest(req2)
}

// Do http request
func (sess *Session) Do(method, url, contentType string, body io.Reader) (*Response, error) {
	req, err := NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return sess.DoRequest(req)
}

// DoWithContext http request
func (sess *Session) DoWithContext(ctx context.Context, method, url, contentType string, body io.Reader) (*Response, error) {
	req, err := NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return sess.DoRequestWithContext(ctx, req)
}

// Get send get request
func (sess *Session) Get(url string) (*Response, error) {
	return sess.Do("GET", url, "", nil)
}

// GetWithContext http request
func (sess *Session) GetWithContext(ctx context.Context, url string) (*Response, error) {
	return sess.DoWithContext(ctx, "GET", url, "", nil)
}

// Post send post request
func (sess *Session) Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("POST", url, contentType, body)
}

// PostWithContext send post request
func (sess *Session) PostWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "POST", url, contentType, body)
}

// PostForm post form request
func (sess *Session) PostForm(url string, data url.Values) (resp *Response, err error) {
	return sess.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// PostFormWithContext post form request
func (sess *Session) PostFormWithContext(ctx context.Context, url string, data url.Values) (resp *Response, err error) {
	return sess.PostWithContext(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Put send put request
func (sess *Session) Put(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("PUT", url, contentType, body)
}

// PutWithContext send put request
func (sess *Session) PutWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "PUT", url, contentType, body)
}

// Delete send delete request
func (sess *Session) Delete(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("DELETE", url, contentType, body)
}

// DeleteWithContext send delete request
func (sess *Session) DeleteWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "DELETE", url, contentType, body)
}
