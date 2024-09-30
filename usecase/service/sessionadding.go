package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func TryAddSession(w http.ResponseWriter, newUserInput *BD.UserInput) (string, error) {
	var SID string
	if newUserInput.TypeUser == "applicant" {
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
		log.Println("BD worker cooky added")
		workerBase.Mu.Lock()
		workerBase.Sessions[SID] = userWorker.ID
		workerBase.Mu.Unlock()

	} else if newUserInput.TypeUser == "employer" {
		employerBase := BD.HandlersEmployer
		userEmployer, ok1 := storage.GetEmployerByEmail(&employerBase, newUserInput.Email)
		log.Println(userEmployer, ok1)

		if ok1 != nil {
			return "", fmt.Errorf(`no user`)
		}
		if ok1 == nil && !storage.EqualHashedPasswords(userEmployer.EmployerPassword, newUserInput.Password) {
			return "", fmt.Errorf(`bad pass`)
		}

		SID = storage.RandStringRunes(32)
		log.Println("BD employer cooky added")
		employerBase.Mu.Lock()
		employerBase.Sessions[SID] = userEmployer.ID
		employerBase.Mu.Unlock()
	}
	return SID, nil
}
