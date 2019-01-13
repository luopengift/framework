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
	Exit(context.Context) error
}

// R inplements Runner interface
type R struct{}

// Prepare prepare
func (r *R) Prepare(ctx context.Context) error {
	return nil
}

// Init init
func (r *R) Init(ctx context.Context) error {
	return nil
}

// Main main
func (r *R) Main(ctx context.Context) error {
	return fmt.Errorf("must rewrite Main Func")
}

// Thread thread
func (r *R) Thread(ctx context.Context) error {
	return nil
}

// Exit exit
func (r *R) Exit(ctx context.Context) error {
	return nil
}

var run Runner = &R{}
