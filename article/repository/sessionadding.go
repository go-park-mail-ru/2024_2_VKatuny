package repository

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// AddSession adding session to db
func AddSession(newUserInput *inmemorydb.UserInput) (string, error) {
	var SID string
	if newUserInput.TypeUser == inmemorydb.WORKER {
		workerBase := inmemorydb.HandlersWorker
		userWorker, ok := GetWorkerByEmail(newUserInput.Email)
		log.Println(userWorker, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if !service.EqualHashedPasswords(userWorker.WorkerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}
		SID = service.RandStringRunes(32)
		log.Println("inmemorydb worker cookie added")
		workerBase.Mu.Lock()
		workerBase.Sessions[SID] = userWorker.ID
		workerBase.Mu.Unlock()

	} else if newUserInput.TypeUser == inmemorydb.EMPLOYER {
		employerBase := inmemorydb.HandlersEmployer
		userEmployer, ok := GetEmployerByEmail(newUserInput.Email)
		log.Println(userEmployer, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if !service.EqualHashedPasswords(userEmployer.EmployerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}

		SID = service.RandStringRunes(32)
		log.Println("inmemorydb employer cookie added")
		employerBase.Mu.Lock()
		employerBase.Sessions[SID] = userEmployer.ID
		employerBase.Mu.Unlock()
	}
	return SID, nil
}
