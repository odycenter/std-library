package errors

import "std-library/app/web/http"

func Forbidden(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "FORBIDDEN"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.Forbidden,
	})
}

func ForbiddenError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "FORBIDDEN",
		severity:   "WARN",
		httpStatus: http.Forbidden,
	})
}

func CustomForbiddenError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.Forbidden,
	})
}
