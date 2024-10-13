package repository

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
)

func TryDellSession(session *http.Cookie) error {
	workerBase := BD.HandlersWorker
	employerBase := BD.HandlersEmployer

	_, ok := GetWorkerBySession(session)

	_, ok1 := GetEmployerBySession(session)

	log.Println(ok, ok1)
	if ok != nil && ok1 != nil {
		return fmt.Errorf(`no sess`)
	}
	if ok == nil {
		log.Println("worker sesion dell")
		workerBase.Mu.Lock()
		delete(workerBase.Sessions, session.Value)
		workerBase.Mu.Unlock()
	} else {
		log.Println("employer sesion dell")
		employerBase.Mu.Lock()
		delete(employerBase.Sessions, session.Value)
		employerBase.Mu.Unlock()
	}
	return nil
}