package commonerrors

import "fmt"

// DB errors
var (
	DBerr = fmt.Errorf("db err")
)

var (
	ErrUnableToCast = fmt.Errorf("unable to cast type")
	ErrInvalidJSON  = fmt.Errorf("invalid JSON, can't unmarshal")
	ErrUnauthorized = fmt.Errorf("unauthorized")
)

// Session errors
var (
	ErrSessionAlreadyExists = fmt.Errorf("session already exists")
	ErrSessionNotFound      = fmt.Errorf("session not found")
)

// Errors to front
var (
	ErrFrontBadSlug          = fmt.Errorf("bad slug")
	ErrFrontUnableToCastSlug = fmt.Errorf("unable to cast slug")
	ErrFrontMethodNotAllowed = fmt.Errorf("method not allowed")
	ErrFrontServiceNotFound  = fmt.Errorf("service not found")
)
