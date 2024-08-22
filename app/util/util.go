package util

import (
	"context"
	internalDB "std-library/app/internal/db"
	internal "std-library/app/internal/module"
	internalscheduler "std-library/app/internal/scheduler"
	internalhttp "std-library/app/internal/web/http"
	"std-library/app/log/consts/logKey"
	"std-library/app/scheduler"
	"time"
)

func NewContext(ctx context.Context) context.Context {
	innerCtx := context.Background()
	if actionLogId := ctx.Value(logKey.Id); actionLogId != nil {
		innerCtx = context.WithValue(ctx, logKey.Id, actionLogId)
	}
	if trace := ctx.Value(logKey.Trace); trace != nil {
		innerCtx = context.WithValue(ctx, logKey.Trace, trace)
	}
	return innerCtx
}

func NewTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(NewContext(ctx), timeout)
}

func NewCancelContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(NewContext(ctx))
}

func NewScheduler() scheduler.Scheduler {
	return internalscheduler.New()
}

func ReadinessProbe(hostURI string) {
	hostname := internal.Hostname(hostURI)
	internal.ResolveHost(context.Background(), hostname)
}

func NewAccessControl() internalhttp.IPv4AccessControl {
	return internalhttp.IPv4AccessControl{}
}

func RegisterOrmLogger() {
	internalDB.ConfigureLog()
}
