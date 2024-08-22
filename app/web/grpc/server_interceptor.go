package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"reflect"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"strings"
	"time"
)

var healthCheckPath = "/grpc.health.v1.Health/Check"

var CustomServerRequestBody func(req interface{}) interface{}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if info.FullMethod == healthCheckPath {
		return handler(ctx, req)
	}

	actionLog := actionlog.Begin(info.FullMethod, "grpc-server")

	contextMap := make(map[string][]any)
	statMap := make(map[string]float64)
	defer func() {
		if er := recover(); er != nil {
			actionLog.AddStat(statMap)

			actionlog.HandleRecover(er, actionLog, contextMap)
			if r, ok := er.(error); ok {
				err = r
			}
		}
	}()

	md, _ := metadata.FromIncomingContext(ctx)

	actionLog.PutContext("controller", reflect.TypeOf(info.Server).String())

	if CustomServerRequestBody != nil {
		actionLog.RequestBody = CustomServerRequestBody(req)
	} else {
		actionLog.RequestBody = req
	}

	if value := md.Get(logKey.RefId); len(value) > 0 {
		actionLog.RefId = value[0]
		ctx = context.WithValue(ctx, logKey.RefId, actionLog.RefId)
	}
	if value := md.Get(logKey.Client); len(value) > 0 {
		if strings.Index(value[0], logKey.ClientPrefix) == 0 {
			if clientValue, err := base64.URLEncoding.DecodeString(value[0][len(logKey.ClientPrefix):]); err == nil {
				actionLog.Client = string(clientValue)
			}
		} else {
			actionLog.Client = value[0]
		}
	}

	if value := md.Get(logKey.ClientHostname); len(value) > 0 {
		actionLog.PutContext(logKey.ClientHostname, value[0])
	}

	if timeout := getServerTimeout(ctx, md); timeout > 0 {
		actionLog.PutContext(timeoutOfDuration, timeout.String())

		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	ctx = context.WithValue(ctx, logKey.Id, actionLog.Id)
	ctx = context.WithValue(ctx, logKey.Action, actionLog.Action)
	ctx = context.WithValue(ctx, logKey.Stat, statMap)
	ctx = context.WithValue(ctx, logKey.Context, contextMap)
	resp, err = handler(ctx, req)
	actionLog.AddContext(contextMap)
	actionLog.AddStat(statMap)

	if err != nil {
		actionlog.HandleRecover(err, actionLog, contextMap)
	} else {
		actionlog.End(actionLog, "ok")
	}

	return
}

func getServerTimeout(ctx context.Context, md metadata.MD) time.Duration {
	if timeout := md.Get(timeoutOfDuration); len(timeout) > 0 {
		d, err := time.ParseDuration(timeout[0])
		if err == nil {
			return d
		}
		slog.ErrorContext(ctx, fmt.Sprintf("parse timeout error:%v", err))
	}

	if enableDefaultTimeout && defaultTimeout > 0 {
		return defaultTimeout
	}

	return 0
}
