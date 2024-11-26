// Package delivery is a handlers layer of session
package delivery

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	"github.com/sirupsen/logrus"
)

type SessionHandlers struct {
	logger         *logrus.Entry
	backendURL     string
	sessionUsecase session.ISessionUsecase
}

func NewSessionHandlers(app *internal.App) *SessionHandlers {
	app.Logger.Debug("Session handlers initialized")
	return &SessionHandlers{
		logger:         &logrus.Entry{Logger: app.Logger},
		backendURL:     app.BackendAddress,
		sessionUsecase: app.Usecases.SessionUsecase,
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

	userType, err := utils.CheckToken(session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgBadUserType,
		})
		return
	}

	userID, err := h.sessionUsecase.CheckAuthorization(r.Context(), userType, session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgNoUserWithSession,
		})
		return
	}
	h.logger.Debugf("%s: got userID: %d", fn, userID)

	user := &dto.JSONUser{
		ID:       userID,
		UserType: userType,
	}
	h.logger.Debugf("%s: user: %v", fn, user)
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

	userWithSession, err := h.sessionUsecase.Login(r.Context(), loginForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      err.Error(), // TODO: formalize error
		})
		return
	}
	h.logger.Debugf("%s: user login successful: %v", fn, userWithSession)

	cookie := utils.MakeAuthCookie(userWithSession.SessionID, h.backendURL)
	http.SetCookie(w, cookie)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body: &dto.JSONUser{
			ID:       userWithSession.ID,
			UserType: userWithSession.UserType,
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

	user, err := h.sessionUsecase.Logout(r.Context(), userType, session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgNoUserWithSession,
		})
		return
	}
	h.logger.Debugf("%s: removed from session and got user: %v", fn, user)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	h.logger.Debugf("%s: deleted user: %v", fn, user)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       user,
	})
}
