package errors

import (
	"path"
	"runtime"
	"strconv"
)

type Code interface {
	ErrorInfo() Info
	ErrorCode() string
	Error() string
	Severity() string
	HTTPStatus() int
}

type Info struct {
	CallerInfo string
	message    string
	Code       int
}

func GetInfo(code int, message ...string) Info {
	return GetInfoBySkip(4, code, message...)
}

func GetInfoBySkip(skip, code int, message ...string) Info {
	_, f, line, ok := runtime.Caller(skip)
	_, file := path.Split(f)
	var result = Info{
		Code: code,
	}
	if ok {
		result.CallerInfo = file + ":" + strconv.Itoa(line)
	}
	if message != nil && len(message) > 0 {
		result.message = message[0]
	} else {
		result.message = strconv.Itoa(code)
	}

	return result
}
