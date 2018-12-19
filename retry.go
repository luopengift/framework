package framework

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/luopengift/log"
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
	defer resp.Body.Close()
	log.Info(string(body))
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

// Retry2 retry
func Retry2(fun func() error, times int) error {

	return nil
}

// RetryByCtl 重试,通过chan来控制
func RetryByCtl(fun func() error, ctl func() <-chan struct{}) error {
	for {
		select {
		case _, ok := <-ctl():
			if !ok {
				return fmt.Errorf("cancelled")
			}
			if err := fun(); err == nil {
				return nil
			}
		}
	}
}
