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

// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string    true         "User's email"
// @Param       password body string    true         "User's password"
// @Success     200      {object}       BD.UserInput
// @Failure     400      {object}       nil
// @Router      /registration/employer/ [post]
func CreateEmployerHandler(h *BD.EmployerHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		storage.SetSecureHeaders(w)
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.EmployerInput)
		decErr := decoder.Decode(newUserInput)
		if decErr != nil {
			w.WriteHeader(400)
			log.Printf("error while unmarshalling employer  JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		if len(newUserInput.EmployerName) < 3 || len(newUserInput.EmployerLastName) < 3 || len(newUserInput.EmployerPosition) < 3 ||
			len(newUserInput.CompanyName) < 3 || len(newUserInput.CompanyDescription) < 10 ||
			len(newUserInput.Website) < 5 || strings.Index(newUserInput.EmployerEmail, "@") < 0 ||
			len(newUserInput.EmployerPassword) < 4 {
			w.WriteHeader(400)
			log.Printf("error while unmarshalling employer  JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		user, err := service.TryCreateEmployer(h, newUserInput)
		if err == nil {
			// UserInputForToken := &BD.UserInput{
			// 	Email:    newUserInput.EmployerEmail,
			// 	Password: newUserInput.EmployerPassword,
			// }
			// LoginFromAnyware(w, UserInputForToken)

			userdata, _ := json.Marshal(user)
			w.Write([]byte(userdata))
			return

		} else {
			w.WriteHeader(400)
			log.Printf("error user with this email already exists: %s", newUserInput.EmployerEmail)
			w.Write([]byte(`{"userAlreadyExist": true}`))
		}

	}
	return http.HandlerFunc(fn)
}
