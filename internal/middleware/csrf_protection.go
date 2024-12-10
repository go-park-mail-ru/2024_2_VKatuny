package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"github.com/sirupsen/logrus"
)

func CSRFProtection(next dto.HandlerFunc, app *internal.App) dto.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "middleware.CSRFProtection"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("WARNING: Can't get logger from context. Processing without it")
		}
		logger.Debugf("%s: got logger from context", fn)

		session, err := r.Cookie(dto.SessionIDName)
		if err == http.ErrNoCookie || session.Value == "" {
			logger.Errorf("checking session: got err %s", err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),
			})
			return
		}

		// Getting userType from session
		userTypeGot, err := utils.CheckToken(session.Value)
		logger.Debugf("got user type %s with token %s", userTypeGot, session.Value)

		if err != nil {
			logger.Errorf("got err %s", err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),
			})
			return
		}

		grpc_request := &grpc_auth.CheckAuthRequest{
			RequestID: r.Context().Value(dto.RequestIDContextKey).(string),
			Session: &grpc_auth.SessionToken{
				ID: session.Value,
			},
		}
		grpc_response, err := app.Microservices.Auth.CheckAuth(r.Context(), grpc_request)
		if err != nil {
			logger.Errorf("got err %s", err)
			UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
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
		ok, err = cryptToken.Check(
			grpc_response.UserData.ID,
			grpc_response.UserData.UserType,
			session.Value,
			token,
		)
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
