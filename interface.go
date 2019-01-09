package framework

import (
	"context"
)

// Preparer prepare app handler interface
type Preparer interface {
	Prepare(context.Context) error
}

// PreparerFunc prepare func
type PreparerFunc func(context.Context) error

// Prepare implements Prepare interface
func (f PreparerFunc) Prepare(ctx context.Context) error {
	return f(ctx)
}

// Initer init app handler interface
type Initer interface {
	Init(context.Context) error
}

// IniterFunc init func
type IniterFunc func(context.Context) error

// Init implements Init interface
func (f IniterFunc) Init(ctx context.Context) error {
	return f(ctx)
}

// Mainer interface
type Mainer interface {
	Main(context.Context) error
}

// MainerFunc main func
type MainerFunc func(context.Context) error

// Main implements Main interface
func (f MainerFunc) Main(ctx context.Context) error {
	return f(ctx)
}

// Threader interface, 用户自行管理携程运行退出等状态, framework仅调起函数
type Threader interface {
	Thread(context.Context) error
}

// ThreaderFunc thread func
type ThreaderFunc func(context.Context) error

// Thread implements Thread interface
func (f ThreaderFunc) Thread(ctx context.Context) error {
	return f(ctx)
}

// Looper interface, 协程循环由framework管理
type Looper interface {
	Loop(context.Context) (exit bool, err error)
}

// LooperFunc thread loop func
type LooperFunc func(context.Context) (bool, error)

// Loop implements Loop interface
func (f LooperFunc) Loop(ctx context.Context) (bool, error) {
	return f(ctx)
}

// Goroutiner interface, 协程循环由framework管理
type Goroutiner interface {
	Loop(context.Context) (exit bool, err error)
}

// GoroutinerFunc thread loop func
type GoroutinerFunc func(context.Context) (bool, error)

// Loop implements Loop interface
func (f GoroutinerFunc) Loop(ctx context.Context) (bool, error) {
	return f(ctx)
}

// Exiter interface
type Exiter interface {
	Exit(ctx context.Context) error
}

// ExiterFunc func
type ExiterFunc func(context.Context) error

// Exit exit
func (f ExiterFunc) Exit(ctx context.Context) error {
	return f(ctx)
}
