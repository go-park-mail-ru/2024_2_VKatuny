package workerDelivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker/repository"
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
	return middleware.HTTPHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(inmemorydb.WorkerInput)
		err := decoder.Decode(newUserInput)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling worker JSON: %s", err)
			return
		}
		if len(newUserInput.WorkerName) < 3 || len(newUserInput.WorkerLastName) < 3 ||
			strings.Index(newUserInput.WorkerEmail, "@") < 0 || len(newUserInput.WorkerPassword) < 4 {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("Bad parameters of the user's fields %d", http.StatusBadRequest)
			return
		}
		user, err := repository.CreateWorker(h, newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, user)
		} else {
			middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.UserAlreadyExist{true})
			log.Printf("error user with this email already exists: %s", newUserInput.WorkerEmail)
		}

	}))
}
