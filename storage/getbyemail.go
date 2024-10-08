package storage

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func GetWorkerByEmail(table *BD.WorkerHandlers, email string) (BD.Worker, error) {
	fmt.Println("$")
	table.Mu.RLock()
	user, err := table.Users[email]
	table.Mu.RUnlock()
	fmt.Println("$$", err)
	fmt.Println(table)
	if err == true {
		return user, nil //fmt.Errorf("User exist")
	} else {
		return BD.Worker{}, fmt.Errorf("No worker with such email")
	}
}

func GetEmployerByEmail(table *BD.EmployerHandlers, email string) (BD.Employer, error) {
	table.Mu.RLock()
	user, ok := table.Users[email]
	table.Mu.RUnlock()
	if ok == true {
		return user, nil //fmt.Errorf("User exist")
	} else {
		return BD.Employer{}, fmt.Errorf("No employer with such email")
	}
}
