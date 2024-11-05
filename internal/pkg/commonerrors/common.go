package commonerrors

import "fmt"

var (
	DBerr = fmt.Errorf("db err")
)

var (
	ErrUnableToCast = fmt.Errorf("unable to cast type") 
)
