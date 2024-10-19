package repository

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// DellSession delete user's session from db
func DellSession(session *http.Cookie) error {
	workerBase := inmemorydb.HandlersWorker
	employerBase := inmemorydb.HandlersEmployer

	_, ok := GetWorkerBySession(session)

	_, ok1 := GetEmployerBySession(session)

	log.Println(ok, ok1)
	if ok != nil && ok1 != nil {
		return fmt.Errorf(`no session`)
	}
	if ok == nil {
		log.Println("worker session dell")
		workerBase.Mu.Lock()
		delete(workerBase.Sessions, session.Value)
		workerBase.Mu.Unlock()
	} else {
		log.Println("employer session dell")
		employerBase.Mu.Lock()
		delete(employerBase.Sessions, session.Value)
		employerBase.Mu.Unlock()
	}
	return nil
}
