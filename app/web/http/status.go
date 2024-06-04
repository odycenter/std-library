package http

// define the HTTP status codes
const (
	OK        = 200
	Created   = 201
	Accepted  = 202
	NoContent = 204

	BadRequest       = 400
	Unauthorized     = 401
	Forbidden        = 403
	NotFound         = 404
	MethodNotAllowed = 405
	NoneAcceptable   = 406
	Conflict         = 409
	Gone             = 410
	UpgradeRequired  = 426
	TooManyRequests  = 429

	InternalServerError = 500
	BadGateway          = 502
	ServiceUnavailable  = 503
	GatewayTimeout      = 504
)
