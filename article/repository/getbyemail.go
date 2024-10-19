package repository

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// GetWorkerByEmail finding worker in db by email
func GetWorkerByEmail(email string) (inmemorydb.Worker, error) {
	table := inmemorydb.HandlersWorker
	table.Mu.RLock()
	user, ok := table.Users[email]
	table.Mu.RUnlock()
	if ok {
		return user, nil //fmt.Errorf("User exist")
	}
	return inmemorydb.Worker{}, fmt.Errorf("No worker with such email")
}

// GetEmployerByEmail finding employer in db by email
func GetEmployerByEmail(email string) (inmemorydb.Employer, error) {
	table := inmemorydb.HandlersEmployer
	table.Mu.RLock()
	user, ok := table.Users[email]
	table.Mu.RUnlock()
	if ok {
		return user, nil //fmt.Errorf("User exist")
	}
	return inmemorydb.Employer{}, fmt.Errorf("No employer with such email")
}
