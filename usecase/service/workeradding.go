package service

import (
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func TryCreateWorker(h *BD.WorkerHandlers, newUserInput *BD.WorkerInput) (BD.Worker, error) {
	_, rErr := storage.GetWorkerByEmail(h, newUserInput.WorkerEmail)

	if rErr != nil {
		hash := storage.HashPassword(newUserInput.WorkerPassword)
		h.Mu.Lock()
		var id uint64 = h.Amount + 1
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
	} else {
		return BD.Worker{}, rErr
	}
}
