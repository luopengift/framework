package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/luopengift/framework"
)

func main() {
	ctx := context.Background()
	app := framework.New()
	app.MainFunc(Main)
	app.ThreadFunc(reportThread)
	app.Run(ctx)
}

func reportThread(ctx context.Context) error {
	for {
		b, err := json.Marshal(framework.NewReport("client-test"))
		if err != nil {
			return err
		}
		reader := bytes.NewBuffer(b)
		framework.Retry("http://127.0.0.1:3456/report", reader, 1)
		time.Sleep(1 * time.Second)
	}
}

// Main mainLoop
func Main(ctx context.Context) error {
	select {
	case <-ctx.Done():
	}
	return nil
}
