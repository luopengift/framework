package framework

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func httpPost(url string, reader io.Reader) error {
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	fmt.Println(body)
	return nil
}

// Retry retry
func Retry(url string, reader io.Reader, retry int) error {
	var err error
	for i := 0; i < retry; i++ {
		if err = httpPost(url, reader); err == nil {
			return nil
		}
	}
	return err
}
