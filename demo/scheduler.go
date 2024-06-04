package demo

import (
	"context"
	"path/filepath"
	"std-library/app/module"
	"std-library/logs"
	reflects "std-library/reflect"
	"strconv"
	"time"
)

type ScheduleModule struct {
	module.Common
}

func (m *ScheduleModule) Initialize() {
	m.Schedule().SetPanicOnAnyAddError(true)

	_, _ = m.Schedule().AddDisallowConcurrent("@every 2s", &TestJob{})

	_, _ = m.Schedule().AddFuncJobWithName(Wrapper("@every 1s", func(ctx context.Context) {
		logs.Info("every 1s")
		// loop 10 times
	}))
}

func Wrapper(spec string, callback func(ctx context.Context)) (string, string, func(context.Context)) {
	name := reflects.FunctionName(callback)
	name = filepath.Base(name)
	return spec, name, callback
}

type TestJob struct {
}

func (t *TestJob) Execute(ctx context.Context) {
	logs.InfoWithCtx(ctx, "trigger every 2s")
	// loop 10 times
	for i := 0; i < 100; i++ {
		logs.Info("sleep 1 sec, loop: " + strconv.Itoa(i))
		time.Sleep(10 * time.Second)
	}
}
