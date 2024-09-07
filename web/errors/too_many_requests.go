package errors

import (
	error "github.com/odycenter/std-library/app/web/errors"
	"github.com/odycenter/std-library/app/web/http"
)

// Deprecated: Use std-library/app/web/errors.TooManyRequests instead.
func TooManyRequests(message string, errorCode ...string) {
	info := error.GetInfoBySkip(2, -1, message)
	code := "TOO_MANY_REQUESTS"
	if errorCode != nil && len(errorCode) > 0 && len(errorCode[0]) > 0 {
		code = errorCode[0]
	}
	panic(error.New(info, code, "WARN", http.TooManyRequests))
}

// Deprecated: Use std-library/app/web/errors.TooManyRequestsError instead.
func TooManyRequestsError(code int, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, "TOO_MANY_REQUESTS", "WARN", http.TooManyRequests))
}

// Deprecated: Use std-library/app/web/errors.CustomTooManyRequestsError instead.
func CustomTooManyRequestsError(code int, errorCode string, message ...string) {
	info := error.GetInfoBySkip(2, code, message...)
	panic(error.New(info, errorCode, "WARN", http.TooManyRequests))
}
