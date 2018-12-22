package framework

import "context"

// Prepare prepare app handler interface
type Prepare interface {
	Prepare(context.Context) error
}

// PrepareFunc prepare func
type PrepareFunc func(context.Context) error

// Prepare implements Prepare interface
func (f PrepareFunc) Prepare(ctx context.Context) error {
	return f(ctx)
}

// Init init app handler interface
type Init interface {
	Init(context.Context) error
}

// InitFunc init func
type InitFunc func(context.Context) error

// Init implements Init interface
func (f InitFunc) Init(ctx context.Context) error {
	return f(ctx)
}

// Main interface
type Main interface {
	Main(context.Context) error
}

// MainFunc main func
type MainFunc func(context.Context) error

// Main implements Main interface
func (f MainFunc) Main(ctx context.Context) error {
	return f(ctx)
}

// Thread interface, 用户自行管理携程运行退出等状态, framework仅调起函数
type Thread interface {
	Thread(context.Context) error
}

// ThreadFunc thread func
type ThreadFunc func(context.Context) error

// Thread implements Thread interface
func (f ThreadFunc) Thread(ctx context.Context) error {
	return f(ctx)
}

// Loop interface, 协程循环由framework管理
type Loop interface {
	Loop(context.Context) (exit bool, err error)
}

// LoopFunc thread loop func
type LoopFunc func(context.Context) (bool, error)

// Loop implements Loop interface
func (f LoopFunc) Loop(ctx context.Context) (bool, error) {
	return f(ctx)
}

// Exiter interface
type Exiter interface {
	Exit(ctx context.Context) error
}

// ExitFunc func
type ExitFunc func(context.Context) error

// Exit exit
func (f ExitFunc) Exit(ctx context.Context) error {
	return f(ctx)
}
