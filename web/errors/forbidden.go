package errors

import (
	error "std-library/app/web/errors"
	"std-library/app/web/http"
)

// Deprecated: use std-library/app/web/errors.Forbidden instead.
func Forbidden(message string, errorCode ...string) {
	info := error.GetInfoBySkip(2, -1, message)
	code := "FORBIDDEN"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.Forbidden))
}

// Deprecated: use std-library/app/web/errors.ForbiddenError instead.
func ForbiddenError(code int, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, "FORBIDDEN", "WARN", http.Forbidden))
}

// Deprecated: use std-library/app/web/errors.CustomForbiddenError instead.
func CustomForbiddenError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, errorCode, "WARN", http.Forbidden))
}
