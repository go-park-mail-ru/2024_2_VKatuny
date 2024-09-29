package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func TryAddSession(w http.ResponseWriter, newUserInput *BD.UserInput) (string, error) {
	workerBase := BD.HandlersWorker
	employerBase := BD.HandlersEmployer
	userWorker, ok := storage.GetWorkerByEmail(&workerBase, newUserInput.Email)
	userEmployer, ok1 := storage.GetEmployerByEmail(&employerBase, newUserInput.Email)
	log.Println(userWorker, ok)
	log.Println(userEmployer, ok1)

	if ok != nil && ok1 != nil {
		return "", fmt.Errorf(`no user`)
	}
	if (ok == nil && !storage.EqualHashedPasswords(userWorker.WorkerPassword, newUserInput.Password)) || (ok1 == nil && !storage.EqualHashedPasswords(userEmployer.EmployerPassword, newUserInput.Password)) {
		return "", fmt.Errorf(`bad pass`)
	}

	SID := storage.RandStringRunes(32)
	if ok == nil {
		workerBase.Mu.RLock()
		workerBase.Sessions[SID] = userWorker.ID
		workerBase.Mu.RUnlock()
	} else {
		employerBase.Mu.RLock()
		fmt.Println("BD", SID)
		employerBase.Sessions[SID] = userEmployer.ID
		employerBase.Mu.RUnlock()
	}
	return SID, nil
}
