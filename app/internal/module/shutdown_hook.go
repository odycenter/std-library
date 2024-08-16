package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	STAGE_0 = 0 // send shutdown signal, begin to shut down processor for external requests, e.g. http server / kafka listener / scheduler

	STAGE_1 = 1 // await external request processor to stop

	STAGE_2 = 2 // after no more new external requests, shutdown internal executors / background tasks

	STAGE_3 = 3 // await internal executors to stop

	STAGE_4 = 4 // after no any task running, shutdown kafka producer

	STAGE_5 = 5 // release all application defined shutdown hook

	STAGE_6 = 6 // release all resources without dependencies, e.g. db / redis / mongo / search

	STAGE_7 = 7 // shutdown kafka log appender, give more time try to forward all logs

	STAGE_8 = 8 // finally, stop the http server, to make sure it responds to incoming requests during shutdown
)

type Stage []func(ctx context.Context, timeoutInMs int64)

var shutdown int32           // 0 for false, 1 for true
var afterShutdownDelay int32 // 0 for false, 1 for true

func IsShutdown() bool {
	return atomic.LoadInt32(&shutdown) == 1
}

func IsAfterShutdownDelay() bool {
	return atomic.LoadInt32(&afterShutdownDelay) == 1
}

type ShutdownHook struct {
	stages                [STAGE_8 + 1]Stage
	shutdownTimeoutInNano int64
	shutdownDelayInSec    int64
}

func (h *ShutdownHook) Initialize() {
	h.shutdownTimeoutInNano = h.getShutdownTimeoutInNano()
	h.shutdownDelayInSec = h.getShutdownDelayInSec()
}

func (h *ShutdownHook) Add(stage int, f func(ctx context.Context, timeoutInMs int64)) {
	if h.stages[stage] == nil {
		h.stages[stage] = make(Stage, 0)
	}
	h.stages[stage] = append(h.stages[stage], f)
}

// in kube env, once Pod is set to the “Terminating” State,
// api-server remove pod from endpoint, and notify kubelet pod deletion simultaneously
// kube-proxy watches endpoint changes, and modify iptables accordingly
// put delay to make sure kube-proxy update iptables before pod stops serving new requests, to reduce connection errors / 503
// (client still needs retry)
func (h *ShutdownHook) getShutdownDelayInSec() int64 {
	shutdownDelay, ok := os.LookupEnv("SHUTDOWN_DELAY_IN_SEC")
	if ok {
		delay, err := strconv.ParseInt(shutdownDelay, 10, 64)
		if err != nil && delay <= 0 {
			log.Panic("shutdown delay must be greater than 0, delay=" + shutdownDelay)
		}
		return delay
	}
	return -1
}

func (h *ShutdownHook) getShutdownTimeoutInNano() int64 {
	shutdownTimeout, ok := os.LookupEnv("SHUTDOWN_TIMEOUT_IN_SEC")
	if ok {
		timeout, _ := strconv.ParseInt(shutdownTimeout, 10, 64)
		timeout = timeout * 1_000_000_000
		if timeout <= 0 {
			log.Panic("shutdown timeout must be greater than 0, timeout=" + shutdownTimeout)
		}
		return timeout
	}
	return 25_000_000_000 // default kube terminationGracePeriodSeconds is 30s, here give 25s try to stop important processes
}

func (h *ShutdownHook) Run() {
	atomic.StoreInt32(&shutdown, 1) // set shutdown to true

	shutdownDelayInSec := h.shutdownDelayInSec
	if shutdownDelayInSec > 0 {
		slog.Info(fmt.Sprintf("delay %v seconds prior to shutdown", shutdownDelayInSec))
		time.Sleep(time.Duration(shutdownDelayInSec) * time.Second)
	}
	atomic.StoreInt32(&afterShutdownDelay, 1) // set afterShutdownDelay to true

	actionLog := actionlog.Begin("app:stop", "app")
	endTime := time.Now().UnixMilli() + h.shutdownTimeoutInNano/1_000_000
	ctx := context.WithValue(context.Background(), logKey.Id, actionLog.Id)
	h.shutdown(ctx, endTime, STAGE_0, STAGE_6)
	actionlog.End(actionLog, "ok") // end action log before closing kafka log appender

	h.shutdown(ctx, endTime, STAGE_7, STAGE_8)
	slog.Info(fmt.Sprintf("shutdown completed, elapsed=%v", actionLog.Elapsed()))
}

func (h *ShutdownHook) shutdown(ctx context.Context, endTime int64, fromStage int, toStage int) {
	for i := fromStage; i <= toStage; i++ {
		stage := h.stages[i]
		if stage == nil {
			continue
		}
		slog.InfoContext(ctx, fmt.Sprintf("shutdown stage: %v", i))
		for _, shutdown := range stage {
			timeoutInMs := endTime - time.Now().UnixMilli()
			if timeoutInMs < 1000 {
				timeoutInMs = 1000 // give 1s if less 1s is left, usually we put larger terminationGracePeriodSeconds than SHUTDOWN_TIMEOUT_IN_SEC, so there are some room to gracefully shutdown rest resources
			}
			shutdown(ctx, timeoutInMs)
		}
	}
}
