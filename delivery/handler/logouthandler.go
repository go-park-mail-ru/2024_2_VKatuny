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
			w.WriteHeader(http.StatusOK)  // client doesn't have a cookie
			return
		}

		errD := service.TryDellSession(session)
		if errD != nil {
			w.WriteHeader(http.StatusOK) // no user with this session
			http.Error(w, `no sess`, http.StatusUnauthorized)
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}
	return http.HandlerFunc(fn)
}
