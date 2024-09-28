package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func CreateEmployerHandler(h *BD.EmployerHandlers) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
			userdata, _ := json.Marshal(user)
			w.Write([]byte(userdata))
		}

	}
	return http.HandlerFunc(fn)
}
