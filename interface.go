package framework

import (
	"context"
)

// PrepareProvider provider prepare goroutine
type PrepareProvider interface {
	PrepareFunc(context.Context) error
}

// PrepareFunc define prepare func type
type PrepareFunc func(context.Context) error

// PrepareFunc implement PrepareProvider interface
func (f PrepareFunc) PrepareFunc(ctx context.Context) error {
	return f(ctx)
}

// InitProvider provider init goroutine
type InitProvider interface {
	InitFunc(context.Context) error
}

// InitFunc define init func type
type InitFunc func(context.Context) error

// InitFunc implement InifProvider interface
func (f InitFunc) InitFunc(ctx context.Context) error {
	return f(ctx)
}

// MainProvider provider main goroutine
type MainProvider interface {
	MainFunc(context.Context) error
}

// MainFunc define main func type
type MainFunc func(context.Context) error

// MainFunc implement MainProvider interface
func (f MainFunc) MainFunc(ctx context.Context) error {
	return f(ctx)
}

// ExitProvider provider exit goroutine
type ExitProvider interface {
	ExitFunc(context.Context) error
}

// ExitFunc define exit func type
type ExitFunc func(context.Context) error

// ExitFunc implement ExitProvider interface
func (f ExitFunc) ExitFunc(ctx context.Context) error {
	return f(ctx)
}

// ThreadProvider provide
type ThreadProvider interface {
	ThreadFunc(context.Context) error
}

// ThreadFunc define thread func typ
type ThreadFunc func(context.Context) error

// ThreadFunc implement ThreadProvider interface
func (f ThreadFunc) ThreadFunc(ctx context.Context) error {
	return f(ctx)
}

// ThreadWithExitProvider provide thread with exit
type ThreadWithExitProvider interface {
	ThreadWithExitFunc(context.Context) (bool, error)
}

// ThreadWithExitFunc define thread func typ
type ThreadWithExitFunc func(context.Context) (bool, error)

// ThreadWithExitFunc implement ThreadWithExitProvider interface
func (f ThreadWithExitFunc) ThreadWithExitFunc(ctx context.Context) (bool, error) {
	return f(ctx)
}
