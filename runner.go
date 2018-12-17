package framework

import "context"

// Runner interface, TODO
type Runner interface {
	Prepare() error
	Init(*App) error
	Main(ctx context.Context, app *App) error
	Thread(ctx context.Context, app *App) error
	Loop(app *App) (bool, error)
	Exit(ctx context.Context, app *App) error
}

// Run inplements Runner interface
type Run struct {
}

// Prepare prepare
func (r *Run) Prepare() error {
	return nil
}

// Init init
func (r *Run) Init(*App) error {
	return nil
}

// Main main
func (r *Run) Main(ctx context.Context, app *App) error {
	return nil
}

// Thread thread
func (r *Run) Thread(ctx context.Context, app *App) error {
	return nil
}

// Loop loop
func (r *Run) Loop(app *App) (bool, error) {
	return false, nil
}

// Exit exit
func (r *Run) Exit(ctx context.Context, app *App) error {
	return nil
}
