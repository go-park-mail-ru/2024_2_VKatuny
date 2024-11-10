package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

// TODO: make access to repositories from context
func RequireAuthorization(next dto.HandlerFunc, repositories *internal.Repositories, userType string) dto.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("WARNING: Can't get logger from context. Processing without it")
		}
		logger.Debug("got logger from context")

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
		if userTypeGot != userType {
			logger.Errorf("forbidden: got user type %s, expected %s", userTypeGot, userType)
			UniversalMarshal(w, http.StatusForbidden, dto.JSONResponse{
				HTTPStatus: http.StatusForbidden,
				Error:      fmt.Errorf("got user type %s, expected %s", userTypeGot, userType).Error(),
			})
			return
		}
		var userID uint64
		if userType == dto.UserTypeApplicant {
			userID, err = repositories.SessionApplicantRepository.GetUserIdBySession(session.Value)
		} else if userType == dto.UserTypeEmployer {
			userID, err = repositories.SessionEmployerRepository.GetUserIdBySession(session.Value)
		}
		if err != nil {
			logger.Errorf("got err %s", err)
			UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),
			})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, dto.UserContextKey, &dto.SessionUser{
			ID:       userID,
			UserType: userType,
		})
		next(w, r.WithContext(ctx))
	}
}
