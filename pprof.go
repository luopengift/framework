package framework

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

func startPprof(ctx context.Context, path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	cpu, err := os.Create(filepath.Join(path, "cpu.prof"))
	if err != nil {
		return err
	}
	defer cpu.Close()

	if err = pprof.StartCPUProfile(cpu); err != nil {
		return err
	}
	defer pprof.StopCPUProfile()

	mem, err := os.Create(filepath.Join(path, "mem.prof"))
	if err != nil {
		return err
	}
	defer mem.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		return err
	}
	go http.ListenAndServe("localhost:6060", nil)
	select {
	case <-ctx.Done():
	}
	return nil
}
