package module

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	internal "github.com/odycenter/std-library/app/internal/module"
	internalscheduler "github.com/odycenter/std-library/app/internal/scheduler"
	"github.com/odycenter/std-library/app/internal/web/sys"
	"github.com/odycenter/std-library/app/scheduler"
	reflects "github.com/odycenter/std-library/reflect"
)

type SchedulerConfig struct {
	DefaultConfig
	scheduler *internalscheduler.SchedulerImpl
}

func (c *SchedulerConfig) Initialize(moduleContext *Context, name string) {
	c.scheduler = internalscheduler.New()
	c.scheduler.PanicOnAnyAddError(false)
	moduleContext.StartupHook.Add(c.scheduler)
	moduleContext.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		c.scheduler.AwaitTermination(ctx, timeoutInMs)
	})

	controller := internal_sys.NewSchedulerController(c.scheduler)
	web.Handler("/_sys/job", controller)
	web.Handler("/_sys/job/*", controller)
}

func (c *SchedulerConfig) AddFuncJob(spec string, process func(ctx context.Context), panicOnAddError ...bool) (scheduler.JobID, error) {
	return c.scheduler.AddFunc(spec, process, panicOnAddError...)
}

func (c *SchedulerConfig) AddFuncJobWithName(spec, action string, process func(ctx context.Context), panicOnAddError ...bool) (scheduler.JobID, error) {
	return c.scheduler.AddFuncWithName(spec, action, process, panicOnAddError...)
}

func (c *SchedulerConfig) Add(spec string, job scheduler.Job, panicOnAddError ...bool) (scheduler.JobID, error) {
	return c.scheduler.Add(spec, job, panicOnAddError...)
}

func (c *SchedulerConfig) AddJobWithName(spec, action string, job scheduler.Job, panicOnAddError ...bool) (scheduler.JobID, error) {
	return c.scheduler.AddWithName(spec, action, job, panicOnAddError...)
}

func (c *SchedulerConfig) AddDisallowConcurrent(spec string, job scheduler.Job, panicOnAddError ...bool) (scheduler.JobID, error) {
	action := reflects.StructName(job)
	return c.scheduler.AddDisallowConcurrentByName(spec, action, internalscheduler.DisallowConcurrent(job), panicOnAddError...)
}

func (c *SchedulerConfig) SetPanicOnAnyAddError(panicOnAnyAddError bool) {
	c.scheduler.PanicOnAnyAddError(panicOnAnyAddError)
}
