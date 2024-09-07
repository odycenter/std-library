package demo

import (
	"context"
	actionlog "github.com/odycenter/std-library/app/log"
	"github.com/odycenter/std-library/app/module"
	reflects "github.com/odycenter/std-library/reflect"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"
)

type ScheduleModule struct {
	module.Common
}

func (m *ScheduleModule) Initialize() {
	m.Schedule().SetPanicOnAnyAddError(true)

	//_, _ = m.Schedule().AddDisallowConcurrent("@every 3s", &TestJob{})

	_, _ = m.Schedule().AddFuncJobWithName(Wrapper("@every 5s", func(ctx context.Context) {
		slog.Debug("every 1s")
		slog.Info("every 1s")
		slog.Warn("every 1s")
		slog.Error("every 1s", "Password", "123")
		actionlog.Context(&ctx, "Password", "123")
	}))

	//_, _ = m.Schedule().AddFuncJobWithName("@every 10s", "funcNameYouDefined", func(ctx context.Context) {
	//	slog.Info("every 1s")
	//})
}

func Wrapper(spec string, callback func(ctx context.Context)) (string, string, func(context.Context)) {
	name := reflects.FunctionName(callback)
	name = filepath.Base(name)
	return spec, name, callback
}

type TestJob struct {
}

func (t *TestJob) Execute(ctx context.Context) {
	slog.InfoContext(ctx, "trigger every 2s")
	// loop 10 times
	for i := 0; i < 100; i++ {
		slog.Info("sleep 1 sec, loop: " + strconv.Itoa(i))
		time.Sleep(10 * time.Second)
	}
}
