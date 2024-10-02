package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

// Login godoc
// @Summary     Realises authentication
// @Description -
// @Tags        Login
// @Accept      json
// @Param       email    body string  true "User's email"
// @Param       password body string  true "User's password"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     401 {object} map[string]interface{}
// @Router      /login/ [post]
func LoginHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		storage.SetSecureHeaders(w)
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.UserInput)
		decErr := decoder.Decode(newUserInput)
		log.Println(newUserInput, decErr)
		if decErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := LoginFromAnyware(w, newUserInput)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
	return http.HandlerFunc(fn)
}

func LoginFromAnyware(w http.ResponseWriter, newUserInput *BD.UserInput) error {
	SID, err := service.TryAddSession(w, newUserInput)
	if err != nil {
		return fmt.Errorf(`no user`)
	}
	log.Println("Cookie received", SID)

	cookie := &http.Cookie{
		Name:     "session_id1",
		Value:    SID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		//Secure:   true, //ubrat
		SameSite: http.SameSiteStrictMode,
		Domain:   BD.BACKENDIP,
	}
	storage.SetSecureHeaders(w)
	http.SetCookie(w, cookie)
	return nil
}
