package grpc

import "time"

var timeoutOfDuration = "timeout_of_duration"
var enableDefaultTimeout = false
var defaultTimeout = 120 * time.Second
var defaultServerTimeoutShift = 300 * time.Millisecond

func EnableTimeout(timeout time.Duration) {
	enableDefaultTimeout = true
	defaultTimeout = timeout
}

func EnableDefaultTimeout() {
	enableDefaultTimeout = true
}
