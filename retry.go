package framework

import (
	"fmt"
	"io"

	"github.com/luopengift/framework/pkg/requests"
	"github.com/luopengift/log"
)

func httpPost(url string, reader io.Reader, timeout int) error {
	req, err := requests.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	resp, err := requests.New().DoRequest(req)
	if err != nil {
		return err
	}
	txt, err := resp.Text()
	if err != nil {
		return err
	}
	log.Info("%v", txt)
	return nil
}

// Retry retry, TODO
func Retry(url string, reader io.Reader, timeout, retry int) error {
	var err error
	for i := 0; i < retry; i++ {
		if err = httpPost(url, reader, timeout); err == nil {
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
