package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
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
	return HttpHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		newUserInput := new(BD.UserInput)
		decErr := decoder.Decode(newUserInput)
		log.Println(newUserInput, decErr)
		if decErr != nil {
			repository.UniversalMarshal(w, http.StatusBadRequest, nil)
			return
		}
		SID, err := service.AddSession(w, newUserInput)
		if err != nil {
			repository.UniversalMarshal(w, http.StatusBadRequest, nil)
			return
		}
		log.Println("Cookie received")
		cookie := &http.Cookie{
			Name:     "session_id1",
			Value:    SID,
			Expires:  time.Now().Add(10 * time.Hour),
			HttpOnly: true,
			//Secure:   true, //ubrat
			SameSite: http.SameSiteStrictMode,
			Domain:   BD.BACKENDIP,
		}
		http.SetCookie(w, cookie)
		if err != nil {
			repository.UniversalMarshal(w, http.StatusUnauthorized, nil)
		}
	}))
}
