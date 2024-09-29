package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func LoginHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.UserInput)
		decErr := decoder.Decode(newUserInput)
		fmt.Println(newUserInput, decErr)
		if decErr != nil {
			w.WriteHeader(401)
			return
		}
		LoginFromAnyware(w, newUserInput)

	}
	return http.HandlerFunc(fn)
}

func LoginFromAnyware(w http.ResponseWriter, newUserInput *BD.UserInput) {
	SID, err := service.TryAddSession(w, newUserInput)

	if err != nil {
		http.Error(w, `no user`, 404)
		return
	}
	fmt.Println("Cooky", SID)
	del := &http.Cookie{
		Name:     "session_id1",
		Value:    SID,
		Expires:  time.Now().Add(-10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, del)

	cookie := &http.Cookie{
		Name:     "session_id1",
		Value:    SID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	storage.SetSecureHeaders(w)
	http.SetCookie(w, cookie)
	//w.Write([]byte("{allok : true}"))
}
