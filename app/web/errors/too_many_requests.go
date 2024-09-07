package errors

import "github.com/odycenter/std-library/app/web/http"

func TooManyRequests(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "TOO_MANY_REQUESTS"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.TooManyRequests,
	})
}

func TooManyRequestsError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "TOO_MANY_REQUESTS",
		severity:   "WARN",
		httpStatus: http.TooManyRequests,
	})
}

func CustomTooManyRequestsError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.TooManyRequests,
	})
}
