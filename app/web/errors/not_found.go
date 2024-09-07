package errors

import "github.com/odycenter/std-library/app/web/http"

func NotFound(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "NOT_FOUND"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.NotFound,
	})
}

func NotFoundError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "NOT_FOUND",
		severity:   "WARN",
		httpStatus: http.NotFound,
	})
}

func CustomNotFoundError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.NotFound,
	})
}
