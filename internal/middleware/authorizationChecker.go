package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/usecase"
	"github.com/sirupsen/logrus"
)

func RequireAuthorization(next http.Handler, logger *logrus.Logger, sessionApplicant, sessionEmployer session.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn := "middleware.RequireAuthorization"

		cookie, err := r.Cookie("session_id1")
		if err == http.ErrNoCookie {
			logger.Errorf("function %s: got %s", fn, err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),  // ErrNoCookie is returned by Request's Cookie method when a cookie is not found.
			})
			return
		} 

		if _, err := usecase.CheckAuthorization(cookie, sessionApplicant, sessionEmployer); err != nil {
			logger.Errorf("function %s: got %s", fn, err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
