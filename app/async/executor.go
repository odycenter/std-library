package async

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"std-library/logs"
)

type IExecutor interface {
	Submit(ctx *context.Context, action string, process func(ctx context.Context))
	Close()
	Running() int
	Free() int
	Waiting() int
}

type Executor struct {
	Name string
	Pool *ants.Pool
}

func New(name string, Size int) IExecutor {
	pool, err := ants.NewPool(Size)
	if err != nil {
		logs.Error("NewExecutor ants.NewPool error:", err)
		return nil
	}
	return &Executor{
		Name: name,
		Pool: pool,
	}
}

func (e *Executor) SubmitTask(ctx *context.Context, action string, task Task) {
	err := e.Pool.Submit(func() {
		execute(ctx, action, task)
	})
	if err != nil {
		logs.Error("Executor Submit error:", err)
	}
}

func (e *Executor) Submit(ctx *context.Context, action string, process func(ctx context.Context)) {
	task := &internalTask{
		process: process,
	}
	e.SubmitTask(ctx, action, task)
}

func (e *Executor) Close() {
	e.Pool.Release()
}

// Running returns the number of workers currently running.
func (e *Executor) Running() int {
	return e.Pool.Running()
}

// Free returns the number of available workers, -1 indicates this pool is unlimited.
func (e *Executor) Free() int {
	return e.Pool.Free()
}

// Waiting returns the number of tasks waiting to be executed.
func (e *Executor) Waiting() int {
	return e.Pool.Waiting()
}
