package employerDelivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/repository"
)

// CreateEmployerHandler creates employers in db
// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string    true         "User's email"
// @Param       password body string    true         "User's password"
// @Success     200      {object}       inmemorydb.UserInput
// @Failure     400      {object}       nil
// @Router      /registration/employer/ [post]
func CreateEmployerHandler(h *inmemorydb.EmployerHandlers) http.Handler {
	return middleware.HTTPHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(inmemorydb.EmployerInput)
		err := decoder.Decode(newUserInput)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling employer  JSON: %s", err)
			return
		}
		if len(newUserInput.EmployerName) < 3 || len(newUserInput.EmployerLastName) < 3 || len(newUserInput.EmployerPosition) < 3 ||
			len(newUserInput.CompanyName) < 3 || len(newUserInput.CompanyDescription) < 10 ||
			len(newUserInput.Website) < 5 || strings.Index(newUserInput.EmployerEmail, "@") < 0 ||
			len(newUserInput.EmployerPassword) < 4 {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling employer  JSON: %s", err)
			return
		}
		user, err := repository.CreateEmployer(h, newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, user)
			return
		}
		middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.UserAlreadyExist{true})
		log.Printf("error user with this email already exists: %s", newUserInput.EmployerEmail)

	}))
}
