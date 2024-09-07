package errors

import (
	error "github.com/odycenter/std-library/app/web/errors"
	"github.com/odycenter/std-library/app/web/http"
)

// Deprecated: std-library/app/web/errors.BadRequest instead.
func BadRequest(message string, errorCode ...string) {
	info := error.GetInfoBySkip(2, -1, message)
	code := "BAD_REQUEST"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.BadRequest))
}

// Deprecated: std-library/app/web/errors.BadRequestError instead.
func BadRequestError(code int, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, "BAD_REQUEST", "WARN", http.BadRequest))
}

// Deprecated: std-library/app/web/errors.CustomBadRequestError instead.
func CustomBadRequestError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, errorCode, "WARN", http.BadRequest))
}
