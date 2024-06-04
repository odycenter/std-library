package scheduler

import (
	"context"
	"time"
)

type Scheduler interface {
	PanicOnAnyAddError(val bool)
	AddFunc(spec string, process func(ctx context.Context), panicOnAddError ...bool) (JobID, error)
	AddFuncWithName(spec, action string, process func(ctx context.Context), panicOnAddError ...bool) (JobID, error)
	Add(spec string, job Job, panicOnAddError ...bool) (JobID, error)
	AddWithName(spec, action string, job Job, panicOnAddError ...bool) (JobID, error)
	Remove(id JobID)
	Start()
	JobsInfo() []JobInfo
	RunningTasks() int
}

type Job interface {
	Execute(ctx context.Context)
}

type JobID int

type JobInfo struct {
	ID      JobID
	Name    string
	Trigger string
	Next    time.Time
	Prev    time.Time
}
