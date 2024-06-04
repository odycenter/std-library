package errors

import "std-library/app/web/http"

func ServiceUnavailable(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "UNAUTHORIZED"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.ServiceUnavailable,
	})
}

func ServiceUnavailableError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "SERVICE_UNAVAILABLE",
		severity:   "WARN",
		httpStatus: http.ServiceUnavailable,
	})
}

func NewServiceUnavailable(code int, message ...string) *Common {
	info := GetInfo(code, message...)
	return &Common{
		Info:       info,
		errorCode:  "SERVICE_UNAVAILABLE",
		severity:   "WARN",
		httpStatus: http.ServiceUnavailable,
	}
}
