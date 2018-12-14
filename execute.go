package framework

import "context"

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
