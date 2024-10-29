// Package sessionDelivery is handlers layer
package sessionDelivery

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
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session/usecase"
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

		if session == nil {
			logger.Errorf("client doesn't have a cookie")
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JsonResponse{
				HttpStatus: http.StatusUnauthorized,
				Error:      "client doesn't have a cookie",
			})
			return
		}

		var id uint64
		var userType string

		if sessionId := session.Value; sessionId != "" {
			logger.WithField("session_id", sessionId).Debug("got session id")

			id, err = repoApplicant.GetUserIdBySession(sessionId)
			userType = dto.UserTypeApplicant

			if err != nil {
				id, err = repoEmployer.GetUserIdBySession(sessionId)
				userType = dto.UserTypeEmployer
			}
		}

		if err == nil {
			logger.Debugf("Just authorized user! UserType: %s; ID %d", userType, id)
			middleware.UniversalMarshal(w, http.StatusOK, dto.JsonResponse{
				HttpStatus: http.StatusOK,
				Body: dto.JsonUserBody{
					UserType: userType,
					ID:       id,
				},
			})
		} else {
			logger.Errorf("authorization error")
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JsonResponse{
				HttpStatus: http.StatusUnauthorized,
				Error:      "authorization error",
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

		newUserInput := new(dto.JsonLoginForm) // for any request and response use DTOs but not a model!
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("can't unmarshal JSON")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}

		sessionId := usecase.GenerateSessionToken()
		if newUserInput.UserType == dto.UserTypeApplicant {
			// need to validate error
			user, _ := repoApplicant.GetByEmail(newUserInput.Email)
			
			err = repoApplicantSession.Add(user.ID, sessionId)
		} else if newUserInput.UserType == dto.UserTypeEmployer {
			// same. Error validation
			user, _ := repoEmployer.GetByEmail(newUserInput.Email)
			err = repoEmployerSession.Add(user.ID, sessionId)
		}

		if err != nil {
			logger.Errorf("err while generating sessionID")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "err while generating session ID",
			})
			return
		}
		logger.Debugf("Cookie received")
		cookie := &http.Cookie{
			Name:     "session_id1", // why id1?
			Value:    sessionId,
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
			middleware.UniversalMarshal(w, http.StatusOK, dto.JsonResponse{
				HttpStatus: http.StatusOK,
				Error:      "client doesn't have a cookie",
			})
			return
		}

		// session is a type of *http.Cookie
		// to delete session i should get session id
		sessionId := session.Value

		err = repoApplicant.Delete(sessionId) // just do it!
		// this is useless because Delete method doesn't return a error if user doesn't exist
		// Oleg it's for you))
		if err != nil {
			err = repoEmployer.Delete(sessionId)
		}

		if err != nil {
			logger.Errorf("no user with this session")
			middleware.UniversalMarshal(w, http.StatusOK, dto.JsonResponse{
				HttpStatus: http.StatusOK,
				Error:      "no user with this session",
			})
			http.Error(w, `no sess`, http.StatusUnauthorized) // hm cheto ne to i vesde tak
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	})
}
