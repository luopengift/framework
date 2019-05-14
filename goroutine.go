package framework

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/luopengift/framework/util"
)

// Goroutine goroutine
type Goroutine struct {
	name     string
	exec     ThreadWithExitProvider
	min, max int
}

func newGoroutine(name string, exec ThreadWithExitProvider, min, max int) *Goroutine {
	return &Goroutine{
		name: name,
		exec: exec,
		min:  min,
		max:  max,
	}
}

// 用户管理goroutine运行退出等状态, framework仅调起函数
func (app *App) runThreads(ctx context.Context) error {
	for id, thread := range app.onThreads {
		if thread == nil {
			return errors.New("thread must not nil")
		}
		wg.Add(1)
		go func(ctx context.Context, id int, execute ThreadProvider) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					app.Log.Fatalf("Thread[%v] %v\n%v", id, err, string(debug.Stack()))
				}
			}()
			if err := execute.ThreadFunc(ctx); err != nil {
				app.Log.Errorf("Thread[%v]: %v", id, err)
			}
		}(ctx, id, thread)
	}
	return nil
}

// SetGoroutineFunc GoroutineFunc
func (app *App) SetGoroutineFunc(name string, fs ThreadWithExitFunc, num ...int) {
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
			return errors.New("goroutine must not nil")
		}
		for i := 0; i < goroutine.min; i++ {
			wg.Add(1)
			go func(ctx context.Context, name string, seq, num int, gor *Goroutine) {
				defer wg.Done()
				var (
					exit  bool // true: 退出goroutine, false: 循环调用goroutine.
					err   error
					entry = func(seq int, execute ThreadWithExitProvider) (bool, error) {
						defer func() {
							if err := recover(); err != nil {
								app.Log.Fatalf("goroutine panic[%v-%v/%v] %v\n%v", name, seq, num, err, string(debug.Stack()))
							}
						}()
						return execute.ThreadWithExitFunc(ctx)
					}
				)

				for !exit {
					select {
					case <-ctx.Done():
						app.Log.Errorf("goroutine ctx[%v-%v/%v]: %v", name, seq, num, ctx.Err())
						return
					default:
						if exit, err = entry(seq, gor.exec); err != nil {
							app.Log.Errorf("goroutine run[%v-%v/%v]: %v", name, seq, num, err)
						}
					}
				}
			}(ctx, goroutine.name, i+1, goroutine.min, goroutine)
		}
	}
	return nil
}
