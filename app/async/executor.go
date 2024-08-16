package async

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"log"
	"log/slog"
	"std-library/app/web/metric"
	"sync"
	"time"
)

type IExecutor interface {
	Submit(ctx *context.Context, action string, process func(ctx context.Context))
	Close()
	Running() int
	Free() int
	Waiting() int
}

var executors sync.Map
var once sync.Once

type Executor struct {
	Name string
	Pool *ants.Pool
}

func New(name string, Size int) IExecutor {
	pool, err := ants.NewPool(Size)
	if err != nil {
		slog.Error("NewExecutor ants.NewPool", "error", err)
		return nil
	}
	_, ok := executors.Load(name)
	if ok {
		log.Panic("Executor already exists:", name)
	}
	executor := &Executor{
		Name: name,
		Pool: pool,
	}
	executors.Store(name, executor)
	once.Do(func() {
		ctx := context.Background()
		go updateMetrics(ctx)
	})
	return executor
}

func (e *Executor) SubmitTask(ctx *context.Context, action string, task Task) {
	err := e.Pool.Submit(func() {
		execute(ctx, action, task)
	})
	if err != nil {
		slog.Error("Executor Submit", "error", err)
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

func updateMetrics(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			collectExecutorMetrics()
		case <-ctx.Done():
			return
		}
	}
}

func collectExecutorMetrics() {
	executors.Range(func(key, value interface{}) bool {
		executor := value.(*Executor)
		metric.ExecutorRunning.WithLabelValues(executor.Name).Set(float64(executor.Running()))
		metric.ExecutorFree.WithLabelValues(executor.Name).Set(float64(executor.Free()))
		metric.ExecutorWaiting.WithLabelValues(executor.Name).Set(float64(executor.Waiting()))
		return true
	})
}
