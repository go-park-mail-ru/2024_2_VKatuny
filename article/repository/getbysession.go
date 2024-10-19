package repository

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// GetWorkerBySession finding worker in db by session
func GetWorkerBySession(session *http.Cookie) (uint64, error) {
	table := inmemorydb.HandlersWorker
	table.Mu.RLock()
	id, ok := table.Sessions[session.Value]
	table.Mu.RUnlock()
	if ok {
		return id, nil
	}
	return 0, fmt.Errorf("No worker with such session")
}

// GetEmployerBySession finding employer in db by session
func GetEmployerBySession(session *http.Cookie) (uint64, error) {
	table := inmemorydb.HandlersEmployer
	table.Mu.RLock()
	fmt.Println(table.Sessions, session.Value)
	id, ok := table.Sessions[session.Value]
	table.Mu.RUnlock()
	if ok {
		return id, nil
	}
	return 0, fmt.Errorf("No employer with such session")

}
