package framework

import "context"

// Prepare prepare app handler interface
type Prepare interface {
	Prepare() error
}

// PrepareFunc prepare func
type PrepareFunc func() error

// Prepare implements Prepar interface
func (f PrepareFunc) Prepare() error {
	return f()
}

// Init init app handler interface
type Init interface {
	Init(*App) error
}

// InitFunc init func
type InitFunc func(*App) error

// Init implements Init interface
func (f InitFunc) Init(app *App) error {
	return f(app)
}

// Main interface
type Main interface {
	Main(ctx context.Context, app *App) error
}

// MainFunc main func
type MainFunc func(ctx context.Context, app *App) error

// Main implements Main interface
func (f MainFunc) Main(ctx context.Context, app *App) error {
	return f(ctx, app)
}

// Thread interface, 用户自行管理携程运行退出等状态, framework仅调起函数
type Thread interface {
	Thread(ctx context.Context, app *App) error
}

// ThreadFunc thread func
type ThreadFunc func(ctx context.Context, app *App) error

// Thread implements Thread interface
func (f ThreadFunc) Thread(ctx context.Context, app *App) error {
	return f(ctx, app)
}

// Loop interface, 协程循环由framework管理
type Loop interface {
	Loop(app *App) (exit bool, err error)
}

// LoopFunc thread loop func
type LoopFunc func(app *App) (bool, error)

// Loop implements Loop interface
func (f LoopFunc) Loop(app *App) (bool, error) {
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
