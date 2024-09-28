package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func CreateWorkerHandler(h *BD.WorkerHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.WorkerInput)
		decErr := decoder.Decode(newUserInput)
		if decErr != nil {
			w.WriteHeader(400)
			log.Printf("error while unmarshalling JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}

		fmt.Println(newUserInput)

		err := service.TryCreateWorker(h, newUserInput)
		if err != nil {
			w.WriteHeader(400)
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
			w.Write([]byte("{}"))
		} else {
			w.Write([]byte("{allok: true}"))
		}

	}
	return http.HandlerFunc(fn)
}
