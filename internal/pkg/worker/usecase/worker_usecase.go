package workerUsecase

import (
	"fmt"
	"strings"
)

func CreateWorkerInputCheck(Name, LastName, Email, Password string) error {
	if len(Name) > 50 || len(LastName) > 50 ||
		strings.Index(Email, "@") < 0 || len(Password) > 50 {
		return fmt.Errorf("worker's fields aren't valid %s %s %s", Name, LastName, Email)
	}
	return nil
}
