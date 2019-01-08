package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/luopengift/framework"
)

func main() {
	framework.ThreadFunc(func(ctx context.Context) error {
		for {
			b, err := json.Marshal(framework.NewReport())
			if err != nil {
				return err
			}
			reader := bytes.NewBuffer(b)
			framework.Retry("http://127.0.0.1:3456/report", reader, 1, 5)
			time.Sleep(1 * time.Second)
		}
	})
	framework.Run()
}
