package requests

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/luopengift/framework/pkg/encoding/json"
)

// Response wrap std response
type Response struct {
	*http.Response
}

// StdLib return net/http.Response
func (resp *Response) StdLib() *http.Response {
	return resp.Response
}

// Text parse parse to string
func (resp *Response) Text() (string, error) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Download parse response to a file
func (resp *Response) Download(name string) (int64, error) {
	defer resp.Body.Close()
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, resp.Body)
}

// JSON parse response
func (resp *Response) JSON(v interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(v)
	resp.Body.Close()
	return err
}
