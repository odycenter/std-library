package async

import (
	"context"
	"path/filepath"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	reflects "std-library/reflect"
	"sync"
)

func RunTask(ctx *context.Context, task Task, wg ...*sync.WaitGroup) {
	RunTaskWithName(ctx, reflects.StructName(task), task, wg...)
}

func RunTaskWithName(ctx *context.Context, action string, task Task, wg ...*sync.WaitGroup) {
	if wg == nil || len(wg) == 0 {
		go execute(ctx, action, task, "task")
		return
	}

	waitGroup := wg[0]
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		execute(ctx, action, task, "task")
	}(waitGroup)
}

func RunFunc(ctx *context.Context, process func(ctx context.Context), wg ...*sync.WaitGroup) {
	action := reflects.FunctionName(process)
	action = filepath.Base(action)
	RunFuncWithName(ctx, action, process, wg...)
}

func RunFuncWithName(ctx *context.Context, action string, process func(ctx context.Context), wg ...*sync.WaitGroup) {
	task := &internalTask{
		process: process,
	}
	RunTaskWithName(ctx, action, task, wg...)
}

func execute(ctx *context.Context, action string, task Task, executeType ...string) {
	actionType := "executor"
	if executeType != nil && len(executeType) > 0 {
		actionType = executeType[0]
	}
	actionLog := actionlog.Begin(action, actionType)
	actionName := "task:" + action
	if ctx != nil {
		rootAction := actionlog.GetAction(ctx)
		if rootAction != "" {
			actionName = rootAction + ":" + actionName
			actionLog.PutContext("root_action", rootAction)
			actionLog.RefId = actionlog.GetId(ctx)
		}
	}
	actionLog.Action = actionName

	contextMap := make(map[string][]any)
	statMap := make(map[string]float64)
	defer func() {
		if err := recover(); err != nil {
			actionLog.AddStat(statMap)
			actionlog.HandleRecover(err, actionLog, contextMap)
		}
	}()

	innerCtx := context.WithValue(context.Background(), logKey.Id, actionLog.Id)
	innerCtx = context.WithValue(innerCtx, logKey.Action, actionName)
	innerCtx = context.WithValue(innerCtx, logKey.Stat, statMap)
	innerCtx = context.WithValue(innerCtx, logKey.Context, contextMap)
	task.Execute(innerCtx)
	actionLog.AddContext(contextMap)
	actionLog.AddContext(contextMap)

	actionlog.End(actionLog, "ok")
}
