package errors

import (
	error "std-library/app/web/errors"
	"std-library/app/web/http"
)

// Deprecated: Use std-library/app/web/errors.Unauthorized instead.
func Unauthorized(message string, errorCode ...string) {
	info := error.GetInfoBySkip(3, -1, message)
	code := "UNAUTHORIZED"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.Unauthorized))
}

// Deprecated: Use std-library/app/web/errors.UnauthorizedError instead.
func UnauthorizedError(code int, message ...string) {
	info := error.GetInfoBySkip(3, code, message...)
	panic(error.New(info, "UNAUTHORIZED", "WARN", http.Unauthorized))
}

// Deprecated: Use std-library/app/web/errors.CustomUnauthorizedError instead.
func CustomUnauthorizedError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(3, code, message...)
	panic(error.New(info, errorCode, "WARN", http.Unauthorized))
}
