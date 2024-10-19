package repository

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// CreateWorker creates worker in db
func CreateWorker(h *inmemorydb.WorkerHandlers, newUserInput *inmemorydb.WorkerInput) (inmemorydb.Worker, error) {
	_, err := GetWorkerByEmail(newUserInput.WorkerEmail)
	log.Println("err ", h, err)
	if err == nil {
		return inmemorydb.Worker{}, fmt.Errorf("user exist")
	}
	hash := service.HashPassword(newUserInput.WorkerPassword)
	var id uint64 = h.Amount + 1
	h.Mu.Lock()
	h.Amount++
	h.Users[newUserInput.WorkerEmail] = inmemorydb.Worker{
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
