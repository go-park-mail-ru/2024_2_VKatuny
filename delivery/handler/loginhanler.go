package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func LoginHandler() http.Handler {
	// fn := func(w http.ResponseWriter, r *http.Request) {
	// 	defer r.Body.Close()

	// 	decoder := json.NewDecoder(r.Body)

	// 	newUserInput1 := new(BD.EmployerInput)
	// 	decErr := decoder.Decode(newUserInput1)
	// 	if decErr != nil {
	// 		w.WriteHeader(400)
	// 		log.Printf("error while unmarshalling employer  JSON: %s", decErr)
	// 		w.Write([]byte("{}"))
	// 		return
	// 	}
	// 	user, err := service.TryCreateEmployer(h, newUserInput1)
	// 	if err != nil {
	// 		w.WriteHeader(400)
	// 		log.Printf("error user with this email already exists: %s", newUserInput1.EmployerEmail)
	// 		w.Write([]byte("{}"))
	// 	} else {
	// 		userdata, _ := json.Marshal(user)
	// 		w.Write([]byte(userdata))
	// 	}

	// }
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoder1 := json.NewDecoder(r.Body)

		newUserInput1 := new(BD.UserInput)
		decErr := decoder1.Decode(newUserInput1)
		fmt.Println(newUserInput1, decErr)
		if decErr != nil {
			return
		}
		api := BD.HandlersWorker
		user, ok := api.Users[newUserInput1.Email]
		fmt.Println(user, ok)
		if !ok {
			http.Error(w, `no user`, 404)
			return
		}
		fmt.Println(storage.EqualHashedPasswords(storage.HashPassword("pass"), "pass"))
		fmt.Println(storage.EqualHashedPasswords(user.WorkerPassword, newUserInput1.Password), newUserInput1.Password, user.WorkerPassword)
		if !storage.EqualHashedPasswords(user.WorkerPassword, newUserInput1.Password) {
			http.Error(w, `bad pass`, 400)
			return
		}

		SID := storage.RandStringRunes(32)
		api.Mu.RLock()
		api.Sessions[SID] = user.ID
		api.Mu.RUnlock()
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    SID,
			Expires:  time.Now().Add(10 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
		w.Write([]byte(SID))

	}
	return http.HandlerFunc(fn)
}
