package employerUsecase

import (
	"fmt"
	"strings"
)

func CreateEmployerInputCheck(Name, LastName, Position, CompanyName, Email, Password string) error {
	if len(Name) > 50 || len(LastName) > 50 || len(Position) > 50 ||
		len(CompanyName) > 50 || strings.Index(Email, "@") < 0 || len(Password) > 50 {
		return fmt.Errorf("employer's fields aren't valid %s %s %s %s %s", Name, LastName, Position, CompanyName, Email)
	}
	return nil
}
