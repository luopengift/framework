package framework

import "context"

// Goroutine goroutine
type Goroutine struct {
	name     string
	before   Function
	exec     FunctionWithExit
	after    Function
	min, max int
}

type null struct{}

func (null) Func(context.Context) error {
	return nil
}

func newGoroutine(name string, exec FunctionWithExit, min, max int) *Goroutine {
	return &Goroutine{
		name:   name,
		before: &null{},
		exec:   exec,
		after:  &null{},
		min:    min,
		max:    max,
	}
}

// Task interface
type Task interface {
	Init(context.Context) error
	exec(context.Context) (bool, error)
	BeforeRun(context.Context) error
	AfterRun(context.Context) error
	RunMode()
}
