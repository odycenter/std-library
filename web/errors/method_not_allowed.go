package errors

import (
	error "github.com/odycenter/std-library/app/web/errors"
	"github.com/odycenter/std-library/app/web/http"
)

// Deprecated: Use std-library/app/web/errors.MethodNotAllowed instead.
func MethodNotAllowed(message string, errorCode ...string) {
	info := error.GetInfoBySkip(2, -1, message)
	code := "METHOD_NOT_ALLOWED"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.MethodNotAllowed))
}

// Deprecated: Use std-library/app/web/errors.MethodNotAllowedError instead.
func MethodNotAllowedError(code int, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, "METHOD_NOT_ALLOWED", "WARN", http.MethodNotAllowed))
}

// Deprecated: Use std-library/app/web/errors.CustomMethodNotAllowedError instead.
func CustomMethodNotAllowedError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, errorCode, "WARN", http.MethodNotAllowed))
}
