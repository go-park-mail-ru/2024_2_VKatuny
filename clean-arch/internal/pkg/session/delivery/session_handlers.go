// Package sessionDelivery is handlers layer
package sessionDelivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session"
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
func AuthorizedHandler(repo session.Repository) http.Handler { // just do it!
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		funcName := "AuthorizedHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
		}
		err := fmt.Errorf("no user with session")
		session, err := r.Cookie("session_id1")
		var id uint64
		var userType string

		if err == nil && session != nil {
			id, err = repo.GetWorkerBySession(session) // just do it!

			userType = repo.APPLICANT // just do it!
			if err != nil {
				id, err = repo.GetEmployerBySession(session) // just do it!
				userType = repo.EMPLOYER                     // just do it!
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
func LoginHandler(repo session.Repository, BACKEND_ADDRESS string) http.Handler { // just do it!
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		funcName := "LoginHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(models.LoginForm) // just do it!
		decErr := decoder.Decode(newUserInput)
		if decErr != nil {
			logger.Errorf("can't unmarshal JSON")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}
		SID, err := repo.AddSession(newUserInput) // just do it!
		if err != nil {
			logger.Errorf("err while generating SID")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "err while generating SID",
			})
			return
		}
		logger.Debugf("Cookie received")
		cookie := &http.Cookie{
			Name:     "session_id1",
			Value:    SID,
			Expires:  time.Now().Add(10 * time.Hour),
			HttpOnly: true,
			//Secure:   true, //s etim ne rabotaet i nado li?
			SameSite: http.SameSiteStrictMode,
			Domain:   BACKEND_ADDRESS,
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
func LogoutHandler(repo session.Repository) http.Handler { // just do it!
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

		errD := repo.DellSession(session) // just do it!
		if errD != nil {
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
