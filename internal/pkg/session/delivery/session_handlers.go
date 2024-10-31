// Package delivery is a handlers layer of session
package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker"

	// "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session"
	sessionUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session/usecase"
	"github.com/sirupsen/logrus"
)

// AuthorizedHandler checks authorization of user
// Authorized godoc
// @Summary     Checks user's authorization
// @Description Gets cookie from user and checks authentication
// @Tags        AuthStatus
// @Param       session_id header string true "Session ID (Cookie)"
// @Success     200
// @Failure     401
// @Router      /authorized [post]
func AuthorizedHandler(repoApplicant session.Repository, repoEmployer session.Repository) http.Handler { // just do it!
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "AuthorizedHandler"

		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
		}

		session, err := r.Cookie("session_id1")

		if err != nil || session == nil {
			logger.Debugf("function: %s; problems with reading cookie", funcName)
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      "client doesn't have a cookie",
			})
		}

		id, userType, sessionID, err := sessionUsecase.SessionCheck(session, repoApplicant, repoEmployer)

		if err == nil {
			logger.WithField("session_id", sessionID).Debug("got session id")
			logger.Debugf("Just authorized user! UserType: %s; ID %d", userType, id)
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Body: dto.JSONUserBody{
					UserType: userType,
					ID:       id,
				},
			})
		} else {
			logger.Errorf("authorization error")
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      err.Error(),
			})
		}
	})
}

// LoginHandler set cookies for users after login
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
func LoginHandler(
	repoApplicantSession session.Repository,
	repoEmployerSession session.Repository,
	repoApplicant worker.Repository,
	repoEmployer employer.Repository,
	backendAddress string,
) http.Handler { // just do it!
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		funcName := "LoginHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(dto.JSONLoginForm) // for any request and response use DTOs but not a model!
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("can't unmarshal JSON")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}

		userID, err := sessionUsecase.LoginCheck(newUserInput, repoApplicant, repoEmployer)
		if err != nil {
			logger.Errorf("function %s: %s", funcName, err.Error())
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JSONResponse{
					HTTPStatus: http.StatusBadRequest,
					Error:      err.Error(),
				},
			)
			return
		}

		sessionID := sessionUsecase.GenerateSessionToken()
		if newUserInput.UserType == dto.UserTypeApplicant {
			err = repoApplicantSession.Add(userID, sessionID)
		} else if newUserInput.UserType == dto.UserTypeEmployer {
			err = repoEmployerSession.Add(userID, sessionID)
		}

		// TODO: remake error comparison
		if err != nil {
			logger.Errorf("function %s: %s", funcName, err.Error())
			middleware.UniversalMarshal(
				w,
				http.StatusBadRequest,
				dto.JSONResponse{
					HTTPStatus: http.StatusBadRequest,
					Error:      err.Error(),
				},
			)
			return
		}

		logger.Debugf("Cookie received")
		cookie := &http.Cookie{
			Name:     "session_id1", // why id1?
			Value:    sessionID,
			Expires:  time.Now().Add(10 * time.Hour),
			HttpOnly: true,
			//Secure:   true, // Enable when using HTTPS
			SameSite: http.SameSiteStrictMode,
			Domain:   backendAddress,
		}
		http.SetCookie(w, cookie)
	})
}

// LogoutHandler deletes cookies when user want to logout
// Logout godoc
// @Summary     Realises deauthentication
// @Description -
// @Tags        Logout
// @Param       session_id header string true "Session ID (Cookie)"
// @Success     200
// @Failure     400
// @Failure     401
// @Router      /logout/ [post]
func LogoutHandler(repoApplicant session.Repository, repoEmployer session.Repository) http.Handler { // just do it!
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		funcName := "LogoutHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
		}
		session, err := r.Cookie("session_id1")
		if err == http.ErrNoCookie {
			logger.Errorf("client doesn't have a cookie")
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Error:      "client doesn't have a cookie",
			})
			return
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(dto.JSONLogoutForm) // for any request and response use DTOs but not a model!
		err = decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("can't unmarshal JSON")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}

		sessionID := session.Value
		err = sessionUsecase.LogoutCheck(newUserInput, sessionID, repoApplicant, repoEmployer)

		if err != nil {
			logger.Errorf("no user with this session")
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Error:      "no user with this session",
			})
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	})
}
