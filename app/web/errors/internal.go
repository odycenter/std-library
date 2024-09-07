package errors

import (
	"github.com/odycenter/std-library/app/web/http"
)

func Internal(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "INTERNAL_ERROR"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "ERROR",
		httpStatus: http.InternalServerError,
	})
}

func InternalError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "INTERNAL_ERROR",
		severity:   "ERROR",
		httpStatus: http.InternalServerError,
	})
}

func CustomInternalError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "ERROR",
		httpStatus: http.InternalServerError,
	})
}
