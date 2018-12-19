package framework

import (
	"context"
	"fmt"
)

// Runner interface, TODO
type Runner interface {
	Prepare(context.Context) error
	Init(context.Context) error
	Main(context.Context) error
	Thread(context.Context) error
	Loop(context.Context) (bool, error)
	Exit(context.Context) error
}

// Run inplements Runner interface
type Run struct {
}

// Prepare prepare
func (r *Run) Prepare(ctx context.Context) error {
	return nil
}

// Init init
func (r *Run) Init(ctx context.Context) error {
	return nil
}

// Main main
func (r *Run) Main(ctx context.Context) error {
	return fmt.Errorf("must rewrite Main Func")
}

// Thread thread
func (r *Run) Thread(ctx context.Context) error {
	return nil
}

// Loop loop
func (r *Run) Loop(ctx context.Context) (bool, error) {
	return false, nil
}

// Exit exit
func (r *Run) Exit(ctx context.Context) error {
	return nil
}

var run Runner = &Run{}
