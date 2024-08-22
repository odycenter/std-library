package internal

import (
	"context"
	"log/slog"
	"std-library/app/async"
	"time"
)

type StartupHook struct {
	// startup has 2 stages, initialize() is for client initialization, e.g. kafka client,
	// start() is to start actual process, like scheduler/listener/etc, those processors may depend on client requires initialize()

	Initialize  []async.Task
	start       []async.Task
	StartStage2 []async.Task
}

func (s *StartupHook) Add(task async.Task) {
	s.start = append(s.start, task)
}

func (s *StartupHook) DoInitialize(ctx context.Context) {
	start := time.Now()
	for _, task := range s.Initialize {
		task.Execute(ctx)
	}
	elapsed := time.Since(start)
	slog.InfoContext(ctx, "after DoInitialize", "elapsed", elapsed.Nanoseconds())
	s.Initialize = []async.Task{}
}

func (s *StartupHook) DoStart(ctx context.Context) {
	start := time.Now()
	for _, task := range s.start {
		task.Execute(ctx)
	}
	for _, task := range s.StartStage2 {
		task.Execute(ctx)
	}
	elapsed := time.Since(start)
	slog.InfoContext(ctx, "after DoStart", "elapsed", elapsed.Nanoseconds())
	s.start = []async.Task{}
	s.StartStage2 = []async.Task{}
}
