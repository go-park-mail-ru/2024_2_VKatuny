package handler

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

// Logout godoc
// @Summary     Realises deauthentication
// @Description -
// @Tags        Logout
// @Param       session_id header string true "Session ID (Cookie)"
// @Success     200 
// @Failure     400 
// @Failure     401 
// @Router      /logout/ [post]
func LogoutHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		storage.SetSecureHeaders(w)
		session, err := r.Cookie("session_id1")
		if err == http.ErrNoCookie {
			return
		}

		errD := service.TryDellSession(session)
		if errD != nil {
			http.Error(w, `no sess`, 401)
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}
	return http.HandlerFunc(fn)
}
