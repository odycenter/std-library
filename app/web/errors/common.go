package errors

type Common struct {
	Info
	errorCode  string
	severity   string
	httpStatus int
}

func New(info Info, errorCode, severity string, httpStatus int) *Common {
	return &Common{
		Info:       info,
		errorCode:  errorCode,
		severity:   severity,
		httpStatus: httpStatus,
	}
}

func (c *Common) ErrorInfo() Info {
	return c.Info
}

func (c *Common) ErrorCode() string {
	if c.errorCode != "" {
		return c.errorCode
	}
	return "UNDEFINED"
}

func (c *Common) Error() string {
	if c.Info.message != "" {
		return c.Info.message
	}
	return c.ErrorCode()
}

func (c *Common) Severity() string {
	if c.errorCode != "" {
		return c.errorCode
	}
	return "ERROR"
}

func (c *Common) HTTPStatus() int {
	return c.httpStatus
}
