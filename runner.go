package framework

import "context"

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
	return nil
}

// Thread thread
func (r *Run) Thread(ctx context.Context) error {
	return nil
}

// Loop loop
func (r *Run) Loop() (bool, error) {
	return false, nil
}

// Exit exit
func (r *Run) Exit(ctx context.Context) error {
	return nil
}

var run = &Run{}
