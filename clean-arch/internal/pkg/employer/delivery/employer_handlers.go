package employerDelivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/inmemorydb"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
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
func CreateEmployerHandler(repo employer.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// logger, ok := r.Context().Value(LoggerKey).(logrus.FieldLogger)

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(models.Employer)
		err := decoder.Decode(newUserInput)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling employer  JSON: %s", err)
			return
		}
		if len(newUserInput.Name) < 3 || len(newUserInput.LastName) < 3 || len(newUserInput.Position) < 3 ||
			len(newUserInput.CompanyName) < 3 || strings.Index(newUserInput.Email, "@") < 0 ||
			len(newUserInput.Password) < 4 {
			middleware.UniversalMarshal(w, http.StatusBadRequest, nil)
			log.Printf("error while unmarshalling employer  JSON: %s", err)
			return
		}
		user, err := repo.Create(newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, user)
			return
		}
		// remove inmemorydb
		middleware.UniversalMarshal(w, http.StatusBadRequest, inmemorydb.UserAlreadyExist{true})
		log.Printf("error user with this email already exists: %s", newUserInput.Email)

	})
}
