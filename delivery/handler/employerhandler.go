package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func CreateEmployerHandler(h *BD.EmployerHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(11)
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.EmployerInput)
		decErr := decoder.Decode(newUserInput)
		if decErr != nil {
			w.WriteHeader(400)
			log.Printf("error while unmarshalling employer  JSON: %s", decErr)
			w.Write([]byte("{}"))
			return
		}
		user, err := service.TryCreateEmployer(h, newUserInput)
		if err != nil {
			w.WriteHeader(400)
			log.Printf("error user with this email already exists: %s", newUserInput.EmployerEmail)
			w.Write([]byte("{}"))
		} else {
			UserInputForToken := &BD.UserInput{
				Email:    newUserInput.EmployerEmail,
				Password: newUserInput.EmployerPassword,
			}
			LoginFromAnyware(w, UserInputForToken)
			storage.SetSecureHeaders(w)
			userdata, _ := json.Marshal(user)
			w.Write([]byte(userdata))
		}

	}
	return http.HandlerFunc(fn)
}
