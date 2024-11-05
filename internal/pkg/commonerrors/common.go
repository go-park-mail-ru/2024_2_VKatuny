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
