package errors

import "std-library/app/web/http"

func Unauthorized(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "UNAUTHORIZED"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.Unauthorized,
	})
}

func UnauthorizedError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "UNAUTHORIZED",
		severity:   "WARN",
		httpStatus: http.Unauthorized,
	})
}

func CustomUnauthorizedError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.Unauthorized,
	})
}
