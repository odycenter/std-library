package dto

import (
	"encoding/json"
	app "github.com/odycenter/std-library/app/conf"
	internallog "github.com/odycenter/std-library/app/internal/log"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	"github.com/odycenter/std-library/app/log/util"
	"log/slog"
	"time"
)

type ActionLog struct {
	Id           string
	Timestamp    time.Time
	result       string
	Action       string
	Client       string
	RefId        string
	Context      map[string][]any
	Stat         map[string]float64
	TraceText    []string
	RequestBody  interface{}
	ResponseBody interface{}
	ErrorCode    string
	ErrorMessage string
	StackTrace   string
	elapsed      int64
	Trace        bool
}

func New() ActionLog {
	return ActionLog{
		Context: make(map[string][]any),
		Stat:    make(map[string]float64),
	}
}

func (actionLog *ActionLog) Begin(action string, actionType string) {
	date := time.Now()
	actionLog.Id = util.GetIDGenerator().Next(date)
	actionLog.Timestamp = date
	actionLog.Action = action
	actionLog.Context = make(map[string][]any)
	actionLog.Context["action_type"] = []any{actionType}
	actionLog.Stat = make(map[string]float64)
}

func (actionLog *ActionLog) PutContext(key string, val ...any) {
	if val == nil || len(val) == 0 {
		return
	}

	v, ok := actionLog.Context[key]
	if !ok {
		v = make([]any, 0, len(val))
	}
	for _, value := range val {
		v = append(v, value)
	}
	actionLog.Context[key] = v
}

func (actionLog *ActionLog) PutStat(key string, val float64) {
	actionLog.Stat[key] = val
}

func (actionLog *ActionLog) Elapsed() int64 {
	return time.Since(actionLog.Timestamp).Nanoseconds()
}

func (actionLog *ActionLog) GetElapsed() int64 {
	return actionLog.elapsed
}

func (actionLog *ActionLog) EnableTrace() {
	actionLog.Trace = true
	actionLog.PutContext(logKey.Trace, true)
}

func (actionLog *ActionLog) Result(result string) *ActionLog {
	actionLog.result = result
	return actionLog
}

func (actionLog *ActionLog) End() *ActionLog {
	actionLog.elapsed = actionLog.Elapsed()
	return actionLog
}

func (actionLog *ActionLog) AddContext(contextMap map[string][]any) {
	if contextMap == nil {
		return
	}
	for k, v := range contextMap {
		actionLog.PutContext(k, v...)
	}
}

func (actionLog *ActionLog) AddStat(statMap map[string]float64) {
	if statMap == nil {
		return
	}
	for k, v := range statMap {
		actionLog.PutStat(k, v)
	}
}

func (actionLog *ActionLog) Output() {
	internallog.Logger.Println(actionLog.String())
	actionLog.Context = nil
	actionLog.Stat = nil
}

func (actionLog *ActionLog) String() string {
	actionLogMessage := map[string]any{
		logKey.Id:    actionLog.Id,
		"app":        app.Name,
		"action":     actionLog.Action,
		"@timestamp": actionLog.Timestamp,
		"result":     actionLog.result,
		"elapsed":    actionLog.elapsed,
	}
	if actionLog.RefId != "" {
		actionLogMessage[logKey.RefId] = actionLog.RefId
	}
	if actionLog.Client != "" {
		actionLogMessage[logKey.Client] = actionLog.Client
	}

	if actionLog.RequestBody != nil {
		switch actionLog.RequestBody.(type) {
		case string:
			actionLogMessage["request_body"] = util.Filter(actionLog.RequestBody.(string), internallog.MaskedFields...)
		default:
			requestBody, _ := json.Marshal(actionLog.RequestBody)
			actionLogMessage["request_body"] = util.Filter(string(requestBody), internallog.MaskedFields...)
		}
	}
	if actionLog.ResponseBody != nil {
		switch actionLog.ResponseBody.(type) {
		case string:
			actionLogMessage["response_body"] = util.Filter(actionLog.ResponseBody.(string), internallog.MaskedFields...)
		default:
			responseBody, _ := json.Marshal(actionLog.ResponseBody)
			actionLogMessage["response_body"] = util.Filter(string(responseBody), internallog.MaskedFields...)
		}
	}
	if actionLog.TraceText != nil && len(actionLog.TraceText) > 0 {
		// loop trace text and add to actionLogMessage, combine with line break
		var traceText string
		for _, trace := range actionLog.TraceText {
			traceText += trace + "\n"
		}
		actionLogMessage["trace"] = traceText
	}
	if actionLog.ErrorCode != "" {
		actionLogMessage["error_code"] = actionLog.ErrorCode
	}
	if actionLog.ErrorMessage != "" {
		actionLogMessage["error_message"] = actionLog.ErrorMessage
	}
	if actionLog.StackTrace != "" {
		actionLogMessage["stack_trace"] = actionLog.StackTrace
	}
	for k, v := range actionLog.Context {
		if len(v) == 1 {
			actionLogMessage["context."+k] = v[0]
			continue
		}
		actionLogMessage["context."+k] = v
	}

	for k, v := range actionLog.Stat {
		actionLogMessage["stat."+k] = v
	}

	actionLogByte, e := json.Marshal(actionLogMessage)
	if e != nil {
		slog.Error(e.Error())
		return ""
	}

	actionLogMessage = nil
	return string(actionLogByte)
}
