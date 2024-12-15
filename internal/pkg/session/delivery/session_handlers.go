// Package delivery is a handlers layer of session
package delivery

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"

	"github.com/sirupsen/logrus"
)

type SessionHandlers struct {
	logger         *logrus.Entry
	backendURL     string
	secretCSRF     string
	authClientGRPC auth_grpc.AuthorizationClient
}

func NewSessionHandlers(app *internal.App) *SessionHandlers {
	app.Logger.Debug("Session handlers initialized")
	if app.Microservices.Auth == nil {
		app.Logger.Fatal("Auth microservice is not initialized")
	}
	return &SessionHandlers{
		logger:         &logrus.Entry{Logger: app.Logger},
		backendURL:     app.BackendAddress,
		secretCSRF:     app.CSRFSecret,
		authClientGRPC: app.Microservices.Auth,
	}
}

// IsAuthorized godoc
// @Summary Check if the user is authorized
// @Description Validates the session cookie and user authorization
// @Tags Session
// @Accept json
// @Produce json
// @Success 200 {object} dto.JSONResponse{body=dto.JSONUser}
// @Failure 401 {object} dto.JSONResponse{error=string}
// @Failure 405 {object} dto.JSONResponse{error=string}
// @Router /api/v1/authorized [get]
func (h *SessionHandlers) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.IsAuthorized"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	session, err := r.Cookie(dto.SessionIDName)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnauthorized,
		})
		return
	}
	h.logger.Debugf("%s: got cookie: %s", fn, session.Value)

	_, err = utils.CheckToken(session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgBadUserType,
		})
		return
	}

	requestID, ok := r.Context().Value(dto.RequestIDContextKey).(string)
	if !ok {
		h.logger.Warningf("%s: no request id provided", fn)
	}
	grpc_request := &auth_grpc.CheckAuthRequest{
		RequestID: requestID,
		// TODO
		Session: &auth_grpc.SessionToken{
			ID: session.Value,
		},
	}

	grpc_response, err := h.authClientGRPC.CheckAuth(r.Context(), grpc_request)
	if err != nil {
		h.logger.Errorf("%s: grpc returned err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgNoUserWithSession,  // TODO: implement error
		})
		return
	}
	
	userData := grpc_response.UserData

	h.logger.Debugf("%s: got userID: %d", fn, userData.ID)
	user := &dto.JSONUser{
		ID:       userData.ID,
		UserType: userData.UserType,
	}
	h.logger.Debugf("%s: user: %v", fn, user)

	cryptToken, err := utils.NewCryptToken(h.secretCSRF)
	if err != nil {
		h.logger.Errorf("%s: can't initialize CSRF token generator %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	tokenCSRF, err := cryptToken.Create(user.ID, user.UserType, session.Value) 
	if err != nil {
		h.logger.Errorf("%s: while creating CSRF token got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: CSRF token created: %s", fn, tokenCSRF)
	w.Header().Set("X-CSRF-Token", tokenCSRF)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       user,
	})
}

// @Summary Log in
// @Description Validates the user credentials and returns a session id
// @Tags Session
// @Accept json
// @Produce json
// @Success 200 {object} dto.JSONResponse{body=dto.JSONUser}
// @Failure 401 {object} dto.JSONResponse{error=string}
// @Failure 405 {object} dto.JSONResponse{error=string}
// @Failure 400 {object} dto.JSONResponse{error=string}
// @Router /api/v1/login [post]
func (h *SessionHandlers) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.Login"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	loginForm := new(dto.JSONLoginForm)
	err := json.NewDecoder(r.Body).Decode(loginForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
		return
	}
	h.logger.Debugf("%s: login form parsed: %v", fn, loginForm)

	requestID, ok := r.Context().Value(dto.RequestIDContextKey).(string)
	if !ok {
		h.logger.Warningf("%s: no request id provided", fn)
	}
	grpc_request := &auth_grpc.AuthRequest{
		RequestID: requestID,
		UserType:  loginForm.UserType,
		Email:     loginForm.Email,
		Password:  loginForm.Password,
	}

	grpc_response, err := h.authClientGRPC.AuthUser(r.Context(), grpc_request)
	if err != nil {
		h.logger.Errorf("%s: grpc returned err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgNoUserWithSession,  // TODO: implement error
		})
		return
	}

	h.logger.Debugf("%s: user login successful: %v", fn, grpc_response)

	cookie := utils.MakeAuthCookie(grpc_response.Session.ID, h.backendURL)
	http.SetCookie(w, cookie)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body: &dto.JSONUser{
			ID:       grpc_response.UserData.ID,
			UserType: grpc_response.UserData.UserType,
		},
	})
}

// @Summary Log out
// @Description Deletes the user session and logs them out
// @Tags Session
// @Accept json
// @Produce json
// @Success 200 {object} dto.JSONResponse{body=dto.JSONUser}
// @Failure 405 {object} dto.JSONResponse{error=string}
// @Failure 500 {object} dto.JSONResponse{error=string}
// @Router /api/v1/logout [post]
func (h *SessionHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.Logout"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	session, err := r.Cookie(dto.SessionIDName)
	if err == http.ErrNoCookie {
		h.logger.Errorf("%s: client doesn't have a cookie", fn)
		middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
			HTTPStatus: http.StatusOK,
			Error:      dto.MsgNoCookie,
		})
		return
	}
	h.logger.Debugf("%s: got cookie: %s", fn, session.Value)

	userType, err := utils.CheckToken(session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
			HTTPStatus: http.StatusOK,
			Error:      dto.MsgBadCookie,
		})
		return
	}
	h.logger.Debugf("%s: got user type: %s", fn, userType)

	requestID, ok := r.Context().Value(dto.RequestIDContextKey).(string)
	if !ok {
		h.logger.Warningf("%s: no request id provided", fn)
	}
	grpc_request := &auth_grpc.DeauthRequest{
		RequestID: requestID,
		Session: &auth_grpc.SessionToken{
			ID: session.Value,
		},
	}
	grpc_response, err := h.authClientGRPC.DeauthUser(r.Context(), grpc_request)
	if err != nil {
		h.logger.Errorf("%s: grpc returned err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgNoUserWithSession,  // TODO: implement error
		})
		return
	}

	h.logger.Debugf("%s: removed from session and got user: %v", fn, grpc_response)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	h.logger.Debugf("%s: deleted user", fn)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}
