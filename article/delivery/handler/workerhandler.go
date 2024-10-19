package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
	//"github.com/go-park-mail-ru/2024_2_VKatuny/article/delivery/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/repository"
)

// CreateWorkerHandler creates worker in db
// CreateWorker godoc
// @Summary     Creates a new user as a worker
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} inmemorydb.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/worker/ [post]
func CreateWorkerHandler(h *inmemorydb.WorkerHandlers) http.Handler {
	return HTTPHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(inmemorydb.WorkerInput)
		err := decoder.Decode(newUserInput)
		if err != nil {
			UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling worker JSON: %s", err)
			return
		}
		if len(newUserInput.WorkerName) < 3 || len(newUserInput.WorkerLastName) < 3 ||
			strings.Index(newUserInput.WorkerEmail, "@") < 0 || len(newUserInput.WorkerPassword) < 4 {
			UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("Bad parameters of the user's fields %d", http.StatusBadRequest)
			return
		}
		user, err := repository.CreateWorker(h, newUserInput)
		if err == nil {
			UniversalMarshal(w, http.StatusOK, user)
		} else {
			UniversalMarshal(w, http.StatusBadRequest, inmemorydb.UserAlreadyExist{true})
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
		}

	}))
}
