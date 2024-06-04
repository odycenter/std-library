package errors

import "std-library/app/web/http"

func MethodNotAllowed(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "METHOD_NOT_ALLOWED"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.MethodNotAllowed,
	})
}

func MethodNotAllowedError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "METHOD_NOT_ALLOWED",
		severity:   "WARN",
		httpStatus: http.MethodNotAllowed,
	})
}

func CustomMethodNotAllowedError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.MethodNotAllowed,
	})
}
