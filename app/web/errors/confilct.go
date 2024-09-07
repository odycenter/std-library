package errors

import "github.com/odycenter/std-library/app/web/http"

func Conflict(message string, errorCode ...string) {
	info := GetInfo(-1, message)
	code := "CONFLICT"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(&Common{
		Info:       info,
		errorCode:  code,
		severity:   "WARN",
		httpStatus: http.Conflict,
	})
}

func ConflictError(code int, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  "CONFLICT",
		severity:   "WARN",
		httpStatus: http.Conflict,
	})
}

func CustomConflictError(code int, errorCode string, message ...string) {
	info := GetInfo(code, message...)
	panic(&Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   "WARN",
		httpStatus: http.Conflict,
	})
}
