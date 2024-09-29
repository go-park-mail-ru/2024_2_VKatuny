package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
	"github.com/go-park-mail-ru/2024_2_VKatuny/usecase/service"
)

func LogoutHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			log.Println(err)
			http.Error(w, `no sess`, 401)
			return
		}

		errD := service.TryDellSession(session)
		if errD != nil {
			http.Error(w, `no sess`, 401)
			return
		}
		storage.SetSecureHeaders(w)
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)

	}
	return http.HandlerFunc(fn)
}
