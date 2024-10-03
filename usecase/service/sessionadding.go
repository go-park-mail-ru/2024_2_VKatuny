package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func AddSession(w http.ResponseWriter, newUserInput *BD.UserInput) (string, error) {
	var SID string
	if newUserInput.TypeUser == BD.WORKER {
		workerBase := BD.HandlersWorker
		userWorker, ok := storage.GetWorkerByEmail(&workerBase, newUserInput.Email)
		log.Println(userWorker, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if ok == nil && !storage.EqualHashedPasswords(userWorker.WorkerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}
		SID = storage.RandStringRunes(32)
		log.Println("BD worker cookie added")
		workerBase.Mu.Lock()
		workerBase.Sessions[SID] = userWorker.ID
		workerBase.Mu.Unlock()

	} else if newUserInput.TypeUser == BD.EMPLOYER {
		employerBase := BD.HandlersEmployer
		userEmployer, ok := storage.GetEmployerByEmail(&employerBase, newUserInput.Email)
		log.Println(userEmployer, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if ok == nil && !storage.EqualHashedPasswords(userEmployer.EmployerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}

		SID = storage.RandStringRunes(32)
		log.Println("BD employer cookie added")
		employerBase.Mu.Lock()
		employerBase.Sessions[SID] = userEmployer.ID
		employerBase.Mu.Unlock()
	}
	return SID, nil
}
