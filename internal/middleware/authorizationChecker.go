package middleware

import (
	"net/http"

	// "github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/usecase"
	// "github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

func RequireAuthorization(next http.Handler, logger *logrus.Logger, sessionApplicant repository.SessionRepository, sessionEmployer repository.SessionRepository, newUserInput *dto.JSONLogoutForm) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn := "middleware.RequireAuthorization"

		cookie, err := r.Cookie("session_id1")
		if err == http.ErrNoCookie {
			logger.Errorf("function %s: got %s", fn, err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(), // ErrNoCookie is returned by Request's Cookie method when a cookie is not found.
			})
			return
		}

		if _, err := usecase.CheckAuthorization(newUserInput, cookie, sessionApplicant, sessionEmployer); err != nil {
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

// func RequireAuthorization(next dto.HandlerFunc, logger *logrus.Logger, repositories *internal.Repositories, userType string) dto.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, err := r.Cookie(dto.SessionIDName)
// 		if err == http.ErrNoCookie || session.Value == "" {
// 			logger.Errorf("checking session: got err %s", err)
// 			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
// 				HTTPStatus: http.StatusUnauthorized,
// 				Error:      err.Error(),
// 			})
// 		}
// 		userType, err := utils.CheckToken(session.Value)
// 		if err != nil {
// 			logger.Errorf("got err %s", err)
// 			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
// 				HTTPStatus: http.StatusUnauthorized,
// 				Error:      err.Error(),
// 			})
// 		}

// 		next(w, r)
// 	}
// }
