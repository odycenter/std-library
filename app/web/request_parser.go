package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	internallog "std-library/app/internal/log"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/log/dto"
	"std-library/app/web/http/header"
	"std-library/nets"
	"strings"
	"time"
)

func ParseRequest(r *http.Request) (context.Context, dto.ActionLog) {
	method := strings.ToLower(r.Method)
	action := "api:" + method + ":" + r.URL.EscapedPath()
	log := actionlog.Begin(action, "web")

	log.PutContext("user_agent", r.UserAgent())
	referer := r.Referer()
	if referer != "" {
		log.PutContext("referer", referer)
	}
	log.PutContext("method", method)
	log.PutContext("client_ip", nets.IP(r).String())
	log.PutContext("remote_address", r.RemoteAddr)
	log.PutContext("http_proto", r.Proto)
	log.PutContext("request_url", r.URL.String())
	queryParams := ParseQueryParams(r)
	if queryParams != "" {
		log.PutContext("query_params", queryParams)
	}
	//contentType := r.Header.Get("Content-Type")
	//if strings.ToLower(strings.Split(contentType, ";")[0]) == "application/x-www-form-urlencoded" {
	RequestBody(log, r)
	//}

	requestId := r.Header.Get("x-request-id")
	if requestId != "" {
		log.PutContext("x_request_id", requestId)
	}

	ctx := context.WithValue(r.Context(), logKey.Id, log.Id)
	ctx = context.WithValue(ctx, logKey.Action, log.Action)

	refId := r.Header.Get(header.RefId)
	if refId != "" {
		log.RefId = refId
	}

	client := r.Header.Get(header.Client)
	if client != "" {
		log.Client = client
	}

	trace := r.Header.Get(header.Trace)
	if trace == "true" {
		log.PutContext("headers", ParseHeaders(r))
		log.EnableTrace()
		ctx = context.WithValue(ctx, header.Trace, trace)
	}
	ctx = context.WithValue(ctx, logKey.ActionLog, log)

	return ctx, log
}

func ParseHeaders(request *http.Request) string {
	var builder strings.Builder
	request.Cookies()
	for key, value := range request.Header {
		if key == "Cookie" {
			continue
		}
		if (builder.Len()) > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		if contains(internallog.MaskedFields, key) {
			builder.WriteString("******")
		} else if len(value) == 1 {
			builder.WriteString(value[0])
		} else {
			builder.WriteString(strings.Join(value, ","))
		}
	}
	return builder.String()
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RequestBody(actionLog dto.ActionLog, request *http.Request) {
	stopwatch := time.Now()
	err := request.ParseForm()
	elapsed := time.Since(stopwatch).Nanoseconds()
	actionLog.Stat["parse_form_elapsed"] = float64(elapsed)
	if err != nil {
		slog.Warn(fmt.Sprintf("parse form error", err.Error()))
		return
	}
	result := parseForm(request)
	if result != nil {
		actionLog.RequestBody = result
	}
}

func parseForm(request *http.Request) map[string]interface{} {
	if request.PostForm == nil || len(request.PostForm) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	for key, value := range request.PostForm {
		if len(value) == 1 {
			result[key] = value[0]
		} else {
			result[key] = value
		}
	}
	return result

}

func ParseQueryParams(request *http.Request) string {
	var builder strings.Builder
	for key, value := range request.URL.Query() {
		if (builder.Len()) > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		if len(value) == 1 {
			builder.WriteString(value[0])
		} else {
			builder.WriteString(strings.Join(value, ","))
		}
	}
	return builder.String()
}
