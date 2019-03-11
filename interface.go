package framework

import (
	"context"
)

// Function func
type Function interface {
	Func(context.Context) error
}

// Func func
type Func func(context.Context) error

// Func func
func (f Func) Func(ctx context.Context) error {
	return f(ctx)
}

// FunctionWithExit with error
type FunctionWithExit interface {
	FuncWithExit(context.Context) (bool, error)
}

// FuncWithExit FuncWithExit
type FuncWithExit func(context.Context) (bool, error)

// FuncWithExit func
func (f FuncWithExit) FuncWithExit(ctx context.Context) (bool, error) {
	return f(ctx)
}
