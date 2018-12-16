package framework

import "context"

// Preparer prepare app handler interface
type Preparer interface {
	Prepare() error
}

// PrepareFunc prepare func
type PrepareFunc func() error

// Prepare implements Preparer interface
func (f PrepareFunc) Prepare() error {
	return f()
}

// Initer init app handler interface
type Initer interface {
	Init(*App) error
}

// InitFunc init func
type InitFunc func(*App) error

// Init implements Initer interface
func (f InitFunc) Init(app *App) error {
	return f(app)
}

// Executer execute interface, 用户自行管理携程运行退出等状态, framework仅调起函数
type Executer interface {
	Execute(ctx context.Context, app *App) error
}

// Func func
type Func func(ctx context.Context, app *App) error

// Execute execute
func (f Func) Execute(ctx context.Context, app *App) error {
	return f(ctx, app)
}

// ExecuteLooper 协程循环由framework管理
type ExecuteLooper interface {
	ExecuteLoop(app *App) (bool, error)
}

// FuncLoop func
type FuncLoop func(app *App) (bool, error)

// ExecuteLoop ExecuteLoop
func (f FuncLoop) ExecuteLoop(app *App) (bool, error) {
	return f(app)
}

// Exiter interface
type Exiter interface {
	Exit(ctx context.Context, app *App) error
}

// ExitFunc func
type ExitFunc func(ctx context.Context, app *App) error

// Exit exit
func (f ExitFunc) Exit(ctx context.Context, app *App) error {
	return f(ctx, app)
}
