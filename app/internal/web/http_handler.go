package web

import (
	"context"
	"encoding/json"
	beegoCtx "github.com/beego/beego/v2/server/web/context"
	actionlog "github.com/odycenter/std-library/app/log"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	appWeb "github.com/odycenter/std-library/app/web"
	"github.com/odycenter/std-library/app/web/errors"
	"net/http"
)

type HTTPHandler struct {
	CustomErrorResponseMessage func(code int, message string) map[string]interface{}
	ErrorWithOkStatus          bool
}

type customResponseWriter struct {
	beegoCtx.Response
	status        int
	contentLength int
}

func (w *customResponseWriter) WriteHeader(code int) {
	w.status = code
	w.Response.WriteHeader(code)
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	length, err := w.Response.Write(b)
	w.contentLength = length
	return length, err
}

func (w *customResponseWriter) GetStatus() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *customResponseWriter) GetContentLength() int {
	return w.contentLength
}

func (f *HTTPHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextMap := make(map[string][]any)
		statMap := make(map[string]float64)
		originalCtx, log := appWeb.ParseRequest(r)
		originalCtx = context.WithValue(originalCtx, logKey.Stat, statMap)
		originalCtx = context.WithValue(originalCtx, logKey.Context, contextMap)
		cw := &customResponseWriter{Response: beegoCtx.Response{ResponseWriter: w}}
		w = cw
		defer func() {
			if err := recover(); err != nil {
				_, ok := log.Context["headers"]
				if !ok {
					log.PutContext("headers", appWeb.ParseHeaders(r))
				}
				log.AddStat(statMap)

				errorCode, code, errorMessage, responseStatus := errorResponse(err)
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				responseStatusCode := responseStatus
				if f.ErrorWithOkStatus {
					responseStatusCode = http.StatusOK
				}
				w.WriteHeader(responseStatusCode)
				response := map[string]interface{}{
					"message": errorMessage,
				}
				if f.CustomErrorResponseMessage != nil {
					customMessage := f.CustomErrorResponseMessage(code, errorMessage)
					for key, value := range customMessage {
						response[key] = value
					}
				}
				response[logKey.Id] = log.Id
				response["errorCode"] = errorCode

				jsonResp, _ := json.Marshal(response)
				w.Write(jsonResp)
				log.PutContext("response_body_length", cw.GetContentLength())
				actionlog.HandleRecover(err, log, contextMap)
			}
		}()

		r = r.WithContext(originalCtx)

		next.ServeHTTP(w, r)

		log.PutContext("response_code", cw.GetStatus())
		log.PutContext("response_body_length", cw.GetContentLength())
		log.AddContext(contextMap)
		log.AddStat(statMap)
		actionlog.End(log, "ok")
	})
}

func errorResponse(e interface{}) (errorCode string, code int, errorMessage string, responseCode int) {
	if err, ok := e.(errors.Code); ok {
		errorMessage = err.Error()
		errorCode = err.ErrorCode()
		responseCode = err.HTTPStatus()
		code = err.ErrorInfo().Code
		return
	}

	errorCode = "INTERNAL_ERROR"
	responseCode = 500
	if err, ok := e.(error); ok {
		errorMessage = err.Error()
	}
	return
}
