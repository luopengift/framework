package framework

import (
	"context"
	"runtime/debug"

	"github.com/luopengift/framework/util"
	"github.com/luopengift/log"
)

// Goroutine goroutine
type Goroutine struct {
	name     string
	before   Function
	exec     FunctionWithExit
	after    Function
	min, max int
}

type null struct{}

func (null) Func(context.Context) error {
	return nil
}

func newGoroutine(name string, exec FunctionWithExit, min, max int) *Goroutine {
	return &Goroutine{
		name:   name,
		before: &null{},
		exec:   exec,
		after:  &null{},
		min:    min,
		max:    max,
	}
}

// Thread interface
func (app *App) Thread(threads ...Function) {
	for _, thread := range threads {
		app.onThreads = append(app.onThreads, thread)
	}
}

// ThreadFunc Func init func program global var
func (app *App) ThreadFunc(fs ...Func) {
	for _, f := range fs {
		app.onThreads = append(app.onThreads, f)
	}
}

// 用户管理goroutine运行退出等状态, framework仅调起函数
func (app *App) runThreads(ctx context.Context) error {
	for id, thread := range app.onThreads {
		if thread == nil {
			return log.Errorf("thread must not nil!")
		}
		wg.Add(1)
		go func(ctx context.Context, id int, execute Function) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					log.Fatal("Thread[%v] %v\n%v", id, err, string(debug.Stack()))
				}
			}()
			if err := execute.Func(ctx); err != nil {
				log.Error("Thread[%v]: %v", id, err)
			}
		}(ctx, id, thread)
	}
	return nil
}

// GoroutineFunc GoroutineFunc
func (app *App) GoroutineFunc(name string, fs FunctionWithExit, num ...int) {
	var min, max int
	switch len(num) {
	case 0:
		min, max = 1, 1
	case 1:
		min, max = num[0], num[0]
	case 2:
		min, max = num[0], num[1]
	default:
		min, max = -1, -1
	}
	if name == "" {
		name = util.Random(8)
	}
	app.goroutines = append(app.goroutines, newGoroutine(name, fs, min, max))
}

// 协程循环由framework管理
func (app *App) runGoroutines(ctx context.Context) error {
	for _, goroutine := range app.goroutines {
		if goroutine.exec == nil {
			return log.Errorf("goroutine must not nil!")
		}
		for i := 0; i < goroutine.min; i++ {
			wg.Add(1)
			go func(ctx context.Context, name string, seq, num int, gor *Goroutine) {
				defer wg.Done()
				var (
					exit  bool // true: 退出goroutine, false: 循环调用goroutine.
					err   error
					entry = func(seq int, execute FunctionWithExit) (bool, error) {
						defer func() {
							if err := recover(); err != nil {
								log.Fatal("goroutine panic[%v-%v/%v] %v\n%v", name, seq, num, err, string(debug.Stack()))
							}
						}()
						return execute.FuncWithExit(ctx)
					}
				)

				for !exit {
					select {
					case <-ctx.Done():
						log.Error("goroutine ctx[%v-%v/%v]: %v", name, seq, num, ctx.Err())
						return
					default:
						if err := gor.before.Func(ctx); err != nil {
							log.Error("before error: %v", err)
							break
						}
						if exit, err = entry(seq, gor.exec); err != nil {
							log.Error("goroutine run[%v-%v/%v]: %v", name, seq, num, err)
						}
						if err := gor.before.Func(ctx); err != nil {
							log.Error("after error: %v", err)
						}
					}
				}
			}(ctx, goroutine.name, i+1, goroutine.min, goroutine)
		}
	}
	return nil
}

// Task interface
type Task interface {
	Init(context.Context) error
	exec(context.Context) (bool, error)
	BeforeRun(context.Context) error
	AfterRun(context.Context) error
	RunMode()
}
