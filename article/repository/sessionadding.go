package repository

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
)

func AddSession(w http.ResponseWriter, newUserInput *BD.UserInput) (string, error) {
	var SID string
	if newUserInput.TypeUser == BD.WORKER {
		workerBase := BD.HandlersWorker
		userWorker, ok := GetWorkerByEmail(newUserInput.Email)
		log.Println(userWorker, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if ok == nil && !service.EqualHashedPasswords(userWorker.WorkerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}
		SID = service.RandStringRunes(32)
		log.Println("BD worker cookie added")
		workerBase.Mu.Lock()
		workerBase.Sessions[SID] = userWorker.ID
		workerBase.Mu.Unlock()

	} else if newUserInput.TypeUser == BD.EMPLOYER {
		employerBase := BD.HandlersEmployer
		userEmployer, ok := GetEmployerByEmail(newUserInput.Email)
		log.Println(userEmployer, ok)

		if ok != nil {
			return "", fmt.Errorf(`no user`)
		}
		if !service.EqualHashedPasswords(userEmployer.EmployerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}

		SID = service.RandStringRunes(32)
		log.Println("BD employer cookie added")
		employerBase.Mu.Lock()
		employerBase.Sessions[SID] = userEmployer.ID
		employerBase.Mu.Unlock()
	}
	return SID, nil
}