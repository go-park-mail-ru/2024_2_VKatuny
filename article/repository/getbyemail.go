package repository

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func GetWorkerByEmail(email string) (BD.Worker, error) {
	table := BD.HandlersWorker
	table.Mu.RLock()
	user, err := table.Users[email]
	table.Mu.RUnlock()
	if err == true {
		return user, nil //fmt.Errorf("User exist")
	} else {
		return BD.Worker{}, fmt.Errorf("No worker with such email")
	}
}

func GetEmployerByEmail(email string) (BD.Employer, error) {
	table := BD.HandlersEmployer
	table.Mu.RLock()
	user, ok := table.Users[email]
	table.Mu.RUnlock()
	if ok == true {
		return user, nil //fmt.Errorf("User exist")
	} else {
		return BD.Employer{}, fmt.Errorf("No employer with such email")
	}
}