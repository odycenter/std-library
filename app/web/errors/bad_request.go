package errors

import (
	"github.com/odycenter/std-library/app/web/http"
)

func BadRequest(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "BAD_REQUEST"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.BadRequest,
	})
}

func BadRequestError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "BAD_REQUEST",
		severity:   "WARN",
		httpStatus: http.BadRequest,
	})
}

func CustomBadRequestError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.BadRequest,
	})
}
