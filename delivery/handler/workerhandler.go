package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

// CreateWorker godoc
// @Summary     Creates a new user as a worker
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} BD.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/worker/ [post]
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
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error while unmarshalling worker JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		if len(newUserInput.WorkerName) < 3 || len(newUserInput.WorkerLastName) < 3 ||
			strings.Index(newUserInput.WorkerEmail, "@") < 0 || len(newUserInput.WorkerPassword) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error while unmarshalling employer  JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		user, err := service.TryCreateWorker(h, newUserInput)
		if err == nil {
			// UserInputForToken := &BD.UserInput{
			// 	Email:    newUserInput.WorkerEmail,
			// 	Password: newUserInput.WorkerPassword,
			// }
			// LoginFromAnyware(w, UserInputForToken)

			userdata, _ := json.Marshal(user)
			log.Println(userdata)
			w.Write([]byte(userdata))
		} else {
			log.Println("!!!", err)
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
			w.Write([]byte(`{"userAlreadyExist": true}`))
		}

	}
	return http.HandlerFunc(fn)
}
