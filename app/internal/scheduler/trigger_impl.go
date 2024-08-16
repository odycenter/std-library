package internal_scheduler

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log/slog"
	"path/filepath"
	"sort"
	internal "std-library/app/internal/module"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/scheduler"
	"std-library/app/web/errors"
	reflects "std-library/reflect"
	"sync"
	"sync/atomic"
	"time"
)

type SchedulerImpl struct {
	cron               *cron.Cron
	panicOnAnyAddError bool
	runningTaskCount   int32
	runningTasks       sync.Map // map[action string]bool
	disallowConcurrent sync.Map // map[action string]JobID
	jobInfo            map[string]scheduler.JobInfo
	jobs               map[string]func(ctx context.Context)
}

func New() *SchedulerImpl {
	s := &SchedulerImpl{}
	s.cron = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	s.jobInfo = make(map[string]scheduler.JobInfo)
	s.jobs = make(map[string]func(ctx context.Context))
	return s
}

func (s *SchedulerImpl) PanicOnAnyAddError(val bool) {
	s.panicOnAnyAddError = val
}

func (s *SchedulerImpl) AddFunc(spec string, process func(ctx context.Context), panicOnAddError ...bool) (scheduler.JobID, error) {
	action := reflects.FunctionName(process)
	action = filepath.Base(action)
	return s.add(spec, action, process, panicOnAddError...)
}

func (s *SchedulerImpl) AddFuncWithName(spec, action string, process func(ctx context.Context), panicOnAddError ...bool) (scheduler.JobID, error) {
	return s.add(spec, action, process, panicOnAddError...)
}

func (s *SchedulerImpl) Add(spec string, job scheduler.Job, panicOnAddError ...bool) (scheduler.JobID, error) {
	action := reflects.StructName(job)
	return s.add(spec, action, job.Execute, panicOnAddError...)
}

func (s *SchedulerImpl) AddWithName(spec, action string, job scheduler.Job, panicOnAddError ...bool) (scheduler.JobID, error) {
	return s.add(spec, action, job.Execute, panicOnAddError...)
}

func (s *SchedulerImpl) AddDisallowConcurrentByName(spec, action string, job DisallowConcurrentJob, panicOnAddError ...bool) (scheduler.JobID, error) {
	id, err := s.add(spec, action, job.job.Execute, panicOnAddError...)
	if err == nil {
		s.disallowConcurrent.Store(action, id)
	}
	return id, err
}

func (s *SchedulerImpl) add(spec, action string, process func(ctx context.Context), panicOnAddError ...bool) (scheduler.JobID, error) {
	_, ok := s.jobInfo[action]
	if ok {
		panic("job already exists, name=" + action)
	}
	p := s.create(action, process)
	entryID, err := s.cron.AddFunc(spec, p)
	if err != nil {
		if s.panicOnAnyAddError || len(panicOnAddError) > 0 && panicOnAddError[0] {
			panic(err)
		}
		return 0, err
	}
	info := scheduler.JobInfo{
		ID:      scheduler.JobID(entryID),
		Name:    action,
		Trigger: spec,
	}
	s.jobInfo[action] = info
	s.jobs[action] = process
	slog.Info("Job register successful", "ID", entryID, "name", action, "spec", entryID)
	return scheduler.JobID(entryID), nil
}

func (s *SchedulerImpl) Remove(id scheduler.JobID) {
	s.cron.Remove(cron.EntryID(id))
}

func (s *SchedulerImpl) info(id scheduler.JobID) scheduler.JobInfo {
	e := s.cron.Entry(cron.EntryID(id))
	return scheduler.JobInfo{
		ID:   scheduler.JobID(e.ID),
		Next: e.Next,
		Prev: e.Prev,
	}
}

func (s *SchedulerImpl) TriggerNow(name string, triggerActionId string) {
	process, ok := s.jobs[name]
	if !ok {
		errors.NotFound("job not found, name=" + name)
	}
	p := s.create(name, process, triggerActionId)
	go p()
}

func (s *SchedulerImpl) JobsInfo() []scheduler.JobInfo {
	var result []scheduler.JobInfo
	for _, entry := range s.jobInfo {
		e := s.cron.Entry(cron.EntryID(entry.ID))
		entry.Next = e.Next
		entry.Prev = e.Prev
		result = append(result, entry)
	}

	slog.Info(fmt.Sprintf("Schedule jobs count: %d ", len(result)))
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (s *SchedulerImpl) RunningTasks() int {
	return int(atomic.LoadInt32(&s.runningTaskCount))
}

func (s *SchedulerImpl) Start() {
	s.cron.Start()
}

func (s *SchedulerImpl) Execute(_ context.Context) {
	s.Start()
}

func (s *SchedulerImpl) AwaitTermination(ctx context.Context, timeoutInMs int64) {
	slog.InfoContext(ctx, "shutting down scheduler")

	innerCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInMs)*time.Millisecond)
	defer cancel()

	select {
	case <-innerCtx.Done():
		slog.InfoContext(innerCtx, fmt.Sprintf("[FAILED_TO_STOP] failed to terminate scheduler, due to timeout, canceledTasks=%d", s.RunningTasks()))
	case <-s.cron.Stop().Done():
		slog.InfoContext(innerCtx, "all jobs have completed")
	}
}

func (s *SchedulerImpl) create(action string, process func(ctx context.Context), triggerActionId ...string) func() {
	return func() {
		atomic.AddInt32(&s.runningTaskCount, 1)
		defer atomic.AddInt32(&s.runningTaskCount, -1)

		id, ok := s.disallowConcurrent.Load(action)
		if ok {
			_, running := s.runningTasks.LoadOrStore(action, true)
			if running {
				info := s.info(id.(scheduler.JobID))
				slog.Warn(fmt.Sprintf("reject job due to disallow Concurrent, %s(id:%v) is still running. previous fire:%v, next fire: %v", action, info.ID, info.Prev, info.Next))
				return
			}

			defer s.runningTasks.Delete(action)
		}

		create(action, process, triggerActionId...)
	}
}

func create(action string, process func(ctx context.Context), triggerActionId ...string) {
	actionName := "job:" + action
	if internal.IsShutdown() {
		slog.Info(fmt.Sprintf("reject job due to server is shutting down!! action: %s", actionName))
		return
	}
	actionLog := actionlog.Begin(actionName, "job")
	if triggerActionId != nil && len(triggerActionId) > 0 {
		actionLog.RefId = triggerActionId[0]
	}

	contextMap := make(map[string][]any)
	statMap := make(map[string]float64)
	defer func() {
		if err := recover(); err != nil {
			actionLog.AddStat(statMap)

			actionlog.HandleRecover(err, actionLog, contextMap)
		}
	}()
	actionLog.PutContext("job", action)
	ctx := context.WithValue(context.Background(), logKey.Id, actionLog.Id)
	ctx = context.WithValue(ctx, logKey.Action, actionName)
	ctx = context.WithValue(ctx, logKey.Stat, statMap)
	ctx = context.WithValue(ctx, logKey.Context, contextMap)
	process(ctx)
	actionLog.AddContext(contextMap)
	actionLog.AddStat(statMap)
	actionlog.End(actionLog, "ok")
}
