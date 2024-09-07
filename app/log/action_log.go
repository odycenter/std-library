package log

import (
	"context"
	"encoding/json"
	internallog "github.com/odycenter/std-library/app/internal/log"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	"github.com/odycenter/std-library/app/log/dto"
	"github.com/odycenter/std-library/app/web/errors"
	"reflect"
	"runtime/debug"
)

func AddMaskedField(fieldName ...string) {
	internallog.AddMaskedField(fieldName...)
}

func Begin(action string, actionType string) dto.ActionLog {
	actionLog := dto.New()
	actionLog.Begin(action, actionType)
	return actionLog
}

func End(actionLog dto.ActionLog, result ...string) {
	if result != nil && len(result) > 0 && result[0] != "" {
		actionLog.Result(result[0])
	}
	actionLog.End().Output()
}

func HandleRecover(r interface{}, actionLog dto.ActionLog, contextMap map[string][]any) {
	actionLog.AddContext(contextMap)
	if err, ok := r.(errors.Code); ok {
		info := err.ErrorInfo()
		if info.CallerInfo != "" {
			actionLog.PutContext("error_from", info.CallerInfo)
		}

		actionLog.ErrorMessage = err.Error()
		actionLog.ErrorCode = err.ErrorCode()
		actionLog.PutContext("response_code", err.HTTPStatus())
	} else if err, ok := r.(error); ok {
		actionLog.ErrorMessage = err.Error()
		actionLog.ErrorCode = "INTERNAL_ERROR"
		actionLog.PutContext("response_code", 500)
	}

	actionLog.StackTrace = string(debug.Stack())
	End(actionLog, "error")
}

func Context(ctx *context.Context, key string, value any) {
	if ctx == nil || key == "" || value == nil {
		return
	}
	contextMap := GetContext(ctx)
	v, ok := contextMap[key]
	if !ok {
		v = make([]any, 0)
	}
	if internallog.IsMaskedField(key) {
		v = append(v, "******")
	} else {
		v = append(v, value)
	}
	contextMap[key] = v
	*ctx = context.WithValue(*ctx, logKey.Context, contextMap)
}

func GetContext(ctx *context.Context) map[string][]any {
	val := getFromContext(ctx, logKey.Context)
	result, ok := val.(map[string][]any)
	if ok {
		return result
	}

	return make(map[string][]any)
}

func Stat(ctx *context.Context, key string, value float64) {
	if ctx == nil || key == "" {
		return
	}
	statMap := stat(ctx)
	statMap[key] = value
	*ctx = context.WithValue(*ctx, logKey.Stat, statMap)
}

func GetStat(ctx *context.Context, key string) float64 {
	statMap := stat(ctx)
	v, ok := statMap[key]
	if ok {
		return v
	}
	return 0
}

func stat(ctx *context.Context) map[string]float64 {
	val := getFromContext(ctx, logKey.Stat)
	result, ok := val.(map[string]float64)
	if ok {
		return result
	}

	return make(map[string]float64)
}

func GetAction(ctx *context.Context) string {
	val := getFromContext(ctx, logKey.Action)
	result, _ := val.(string)
	return result
}

func GetId(ctx *context.Context) string {
	val := getFromContext(ctx, logKey.Id)
	result, _ := val.(string)
	return result
}

func getFromContext(ctx *context.Context, key string) any {
	if ctx == nil {
		return nil
	}

	val := (*ctx).Value(key)
	if val == nil {
		return nil
	}

	return val
}

func RequestBody(req interface{}, traceLog bool) map[string]any {
	if !traceLog || reflect.TypeOf(req).String() != "*common.Request" {
		return nil
	}

	requestString, _ := json.Marshal(req)
	var requestMap = make(map[string]interface{})
	if err := json.Unmarshal(requestString, &requestMap); err != nil {
		return nil
	}

	dataString, dataExist := requestMap["Data"].(string)
	if dataExist {
		var dataMap = make(map[string]interface{})
		err2 := json.Unmarshal([]byte(dataString), &dataMap)
		if err2 == nil {
			requestMap["Data"] = dataMap
		}
	}
	return requestMap

}
