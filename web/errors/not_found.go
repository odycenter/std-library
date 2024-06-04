package errors

import (
	error "std-library/app/web/errors"
	"std-library/app/web/http"
)

// Deprecated: Use std-library/app/web/errors.NotFound instead.
func NotFound(message string, errorCode ...string) {
	info := error.GetInfoBySkip(3, -1, message)
	code := "NOT_FOUND"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.NotFound))
}

// Deprecated: Use std-library/app/web/errors.NotFoundError instead.
func NotFoundError(code int, message ...string) {
	info := error.GetInfoBySkip(3, code, message...)
	panic(error.New(info, "NOT_FOUND", "WARN", http.NotFound))
}

// Deprecated: Use std-library/app/web/errors.CustomNotFoundError instead.
func CustomNotFoundError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(3, code, message...)
	panic(error.New(info, errorCode, "WARN", http.NotFound))
}
