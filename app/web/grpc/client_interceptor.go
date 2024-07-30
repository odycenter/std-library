package grpc

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	app "std-library/app/conf"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"strconv"
	"time"
)

var CustomClientRequestBody func(req interface{}) interface{}

var slowOperationThresholdInNanos = 20 * time.Second.Nanoseconds()

func ClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	actionLog := actionlog.Begin(method, "grpc-client")
	if existsId := ctx.Value(logKey.Id); existsId != nil {
		id := existsId.(string)
		actionLog.Id = id
	}
	actionLog.PutContext("conn_target", cc.Target())
	defer func() {
		if err := recover(); err != nil {
			actionlog.HandleRecover(err, actionLog, nil)
		}
	}()

	ctx = metadata.AppendToOutgoingContext(ctx, logKey.RefId, actionLog.Id,
		logKey.Client, logKey.ClientPrefix+base64.URLEncoding.EncodeToString([]byte(app.Name)),
		logKey.ClientHostname, app.LocalHostName())

	traceLog := false
	if trace := ctx.Value(logKey.Trace); trace != nil && trace.(string) == "true" {
		traceLog = true
		if CustomClientRequestBody != nil {
			actionLog.RequestBody = CustomClientRequestBody(req)
		} else {
			actionLog.RequestBody = req
		}
		ctx = metadata.AppendToOutgoingContext(ctx, logKey.Trace, strconv.FormatBool(traceLog))
	}

	var timeout = getClientTimeout(ctx)
	if timeout > 0 {
		actionLog.PutContext(timeoutOfDuration, timeout)
		ctx = metadata.AppendToOutgoingContext(ctx, timeoutOfDuration, timeout.String())
	}

	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		actionlog.HandleRecover(err, actionLog, nil)
	} else {
		actionLog.End()
		if actionLog.Elapsed() > slowOperationThresholdInNanos {
			actionLog.PutContext("slow_grpc", true)
			actionLog.Result("warn")
		} else {
			actionLog.Result("ok")
		}
		actionlog.Output(actionLog)
	}

	return err
}

func getClientTimeout(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if ok {
		timeout := time.Until(deadline) - defaultServerTimeoutShift
		if timeout > time.Second {
			return timeout
		}
		return time.Second
	}
	if enableDefaultTimeout {
		return defaultTimeout
	}

	return 0
}

func SlowOperationThreshold(threshold time.Duration) {
	slowOperationThresholdInNanos = threshold.Nanoseconds()
}
