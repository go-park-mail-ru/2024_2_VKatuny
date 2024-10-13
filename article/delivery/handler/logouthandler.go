package handler

import (
	"net/http"
	"time"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/article/delivery/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/repository"
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
	return HttpHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		session, err := r.Cookie("session_id1")
		if err == http.ErrNoCookie {
			UniversalMarshal(w, http.StatusOK, nil) // client doesn't have a cookie
			return
		}

		errD := repository.TryDellSession(session)
		if errD != nil {
			UniversalMarshal(w, http.StatusOK, nil) // no user with this session
			http.Error(w, `no sess`, http.StatusUnauthorized)
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}))
}
