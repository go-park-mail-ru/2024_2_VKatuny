// Package delivery is a handlers layer of session
package delivery

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	"github.com/sirupsen/logrus"
)

type SessionHandlers struct {
	logger           *logrus.Entry
	backendURL       string
	sessionUsecase   session.ISessionUsecase
	applicantUsecase applicant.IApplicantUsecase
	employerUsecase  employer.IEmployerUsecase
}

func NewSessionHandlers(app *internal.App) *SessionHandlers {
	app.Logger.Debug("Session handlers initialized")
	return &SessionHandlers{
		logger:         &logrus.Entry{Logger: app.Logger},
		sessionUsecase: app.Usecases.SessionUsecase,
		backendURL:     app.BackendAddress,
	}
}

func (h *SessionHandlers) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.IsAuthorized"

	session, err := r.Cookie(dto.SessionIDName)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnauthorized,
		})
		return
	}
	// TODO: session ID should be a field. Usable method in other PR
	h.logger.Debugf("%s: got cookie: %s", fn, session.Value)

	// TODO: remake this. It should be a session usecase not utils
	userType, err := utils.CheckToken(session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgBadUserType,
		})
		return
	}

	userID, err := h.sessionUsecase.CheckAuthorization(userType, session.Value)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgNoUserWithSession,
		})
		return
	}
	// TODO: userID should be a field
	h.logger.Debugf("%s: got userID: %d", fn, userID)

	var user interface{}
	if userType == dto.UserTypeApplicant {
		// TODO: think about usecase that i should use (session's or applicant's)
		user, err = h.applicantUsecase.GetByID(userID)
		if err != nil {
			h.logger.Errorf("%s: got err %s", fn, err)
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      dto.MsgDataBaseError,
			})
			return
		}

	} else if userType == dto.UserTypeEmployer {
		user, err = h.employerUsecase.GetByID(userID)
		if err != nil {
			h.logger.Errorf("%s: got err %s", fn, err)
			middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
				HTTPStatus: http.StatusUnauthorized,
				Error:      dto.MsgDataBaseError,
			})
			return
		}
	}

	h.logger.Debugf("%s: user: %v", fn, user)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       user,
	})
}

func (h *SessionHandlers) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.Login"

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

	userWithSession, err := h.sessionUsecase.Login(loginForm)
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

func (h *SessionHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SessionHandlers.Logout"

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

	user, err := h.sessionUsecase.Logout(userType, session.Value)
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
