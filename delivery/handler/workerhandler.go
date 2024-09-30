package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func CreateWorkerHandler(h *BD.WorkerHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		storage.SetSecureHeaders(w)
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.WorkerInput)
		decErr := decoder.Decode(newUserInput)
		if decErr != nil {
			w.WriteHeader(400)
			log.Printf("error while unmarshalling worker JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		user, err := service.TryCreateWorker(h, newUserInput)
		if err == nil {
			UserInputForToken := &BD.UserInput{
				Email:    newUserInput.WorkerEmail,
				Password: newUserInput.WorkerPassword,
			}
			LoginFromAnyware(w, UserInputForToken)

			userdata, _ := json.Marshal(user)
			log.Println(userdata)
			w.Write([]byte(userdata))
		} else {
			log.Println("!!!", err)
			w.WriteHeader(400)
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
			w.Write([]byte(`{"userAlreadyExist": true}`))
		}

	}
	return http.HandlerFunc(fn)
}
