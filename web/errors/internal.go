package errors

import (
	error "github.com/odycenter/std-library/app/web/errors"
	"github.com/odycenter/std-library/app/web/http"
)

// Deprecated: Use std-library/app/web/errors.Internal instead.
func Internal(message string, errorCode ...string) {
	info := error.GetInfoBySkip(2, -1, message)
	code := "INTERNAL_ERROR"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.InternalServerError))
}

// Deprecated: Use std-library/app/web/errors.InternalError instead.
func InternalError(code int, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, "INTERNAL_ERROR", "WARN", http.InternalServerError))
}

// Deprecated: Use std-library/app/web/errors.CustomInternalError instead.
func CustomInternalError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, errorCode, "ERROR", http.InternalServerError))
}
