package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func CSRFProtection(next HandlerFunc, app *internal.App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "middleware.CSRFProtection"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("WARNING: Can't get logger from context. Processing without it")
		}
		logger.Debugf("%s: got logger from context", fn)

		user, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
		if !ok {
			logger.Errorf("%s: got err %s", fn, "failed to get user from context")
			UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
				Error:      "failed to get user from context",
			})
			return
		}
		logger.Debugf("%s: got user from context: %v", fn, user)

		session, err := r.Cookie(dto.SessionIDName)
		if err == http.ErrNoCookie || session.Value == "" {
			logger.Errorf("checking session: got err %s", err)
			UniversalMarshal(w, http.StatusForbidden, dto.JSONResponse{
				HTTPStatus: http.StatusForbidden,
				Error:      err.Error(),
			})
			return
		}

		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			logger.Errorf("%s: csrf token is empty", fn)
			UniversalMarshal(w, http.StatusForbidden, dto.JSONResponse{
				HTTPStatus: http.StatusForbidden,
				Error:      "csrf token is empty",
			})
			return
		}
		logger.Debugf("%s: got token: %s", fn, token)

		cryptToken := utils.NewCryptToken(app.CSRFSecret)
		ok, err = cryptToken.Check(user.ID, user.UserType, session.Value, token)
		if err != nil {
			errMsg := commonerrors.ErrUncoveredError
			if err == utils.ErrTokenExpired {
				errMsg = commonerrors.ErrFrontCSRFExpired
			}
			logger.Errorf("%s: got err %s", fn, err)
			UniversalMarshal(w, http.StatusForbidden, dto.JSONResponse{
				HTTPStatus: http.StatusForbidden,
				Error:      errMsg.Error(),
			})
			return
		}

		if !ok {
			logger.Errorf("%s: token doesn't match", fn)
			UniversalMarshal(w, http.StatusForbidden, dto.JSONResponse{
				HTTPStatus: http.StatusForbidden,
				Error:      commonerrors.ErrFrontCSRFTokenDoesntMatch.Error(),
			})
			return
		}

		next(w, r)
	}
}
