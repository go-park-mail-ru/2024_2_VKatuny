package repository

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func GetWorkerBySession(session *http.Cookie) (uint64, error) {
	table := BD.HandlersWorker
	table.Mu.RLock()
	id, err := table.Sessions[session.Value]
	table.Mu.RUnlock()
	if err == true {
		return id, nil
	} else {
		return 0, fmt.Errorf("No worker with such session")
	}
}

func GetEmployerBySession(session *http.Cookie) (uint64, error) {
	table := BD.HandlersEmployer
	table.Mu.RLock()
	fmt.Println(table.Sessions, session.Value)
	id, err := table.Sessions[session.Value]
	table.Mu.RUnlock()
	if err == true {
		return id, nil
	} else {
		return 0, fmt.Errorf("No employer with such session")
	}
}
