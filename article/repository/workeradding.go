package repository

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
)

func TryCreateWorker(h *BD.WorkerHandlers, newUserInput *BD.WorkerInput) (BD.Worker, error) {
	_, err := GetWorkerByEmail(newUserInput.WorkerEmail)
	log.Println("err ", h, err)
	if err == nil {
		return BD.Worker{}, fmt.Errorf("User exist")
	} else {
		hash := service.HashPassword(newUserInput.WorkerPassword)
		var id uint64 = h.Amount + 1
		h.Mu.Lock()
		h.Amount += 1
		h.Users[newUserInput.WorkerEmail] = BD.Worker{
			ID:              id,
			WorkerName:      newUserInput.WorkerName,
			WorkerLastName:  newUserInput.WorkerLastName,
			WorkerBirthDate: newUserInput.WorkerBirthDate,
			WorkerEmail:     newUserInput.WorkerEmail,
			WorkerPassword:  hash,
		}
		h.Mu.Unlock()
		log.Println("worker registrated")
		return h.Users[newUserInput.WorkerEmail], nil
	}
}