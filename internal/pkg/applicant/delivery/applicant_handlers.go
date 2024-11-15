package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// CreateWorkerHandler creates applicant in db
// CreateWorker godoc
// @Summary     Creates a new user as a applicant
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} inmemorydb.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/applicant/ [post]
func (h *ApplicantHandlers) ApplicantRegistration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantRegistration"

	applicantRegistrationForm := new(dto.JSONApplicantRegistrationForm)

	err := json.NewDecoder(r.Body).Decode(applicantRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
		return
	}
	h.logger.Debugf("%s: json decoded: %v", fn, applicantRegistrationForm)

	// TODO: add json validation with govalidator

	applicant, err := h.applicantUsecase.Create(applicantRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgUnableToCreateUser,
		})
		return
	}
	h.logger.Debugf("%s: user created successfully: %v", fn, applicant)

	loginForm := &dto.JSONLoginForm{
		UserType: dto.UserTypeApplicant,
		Email:    applicantRegistrationForm.Email,
		Password: applicantRegistrationForm.Password,
	}
	userWithSession, err := h.sessionUsecase.Login(loginForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      err.Error(), // TODO: standardize errors
		})
		return
	}
	h.logger.Debugf("%s: user logged in: %v", fn, userWithSession)

	cookie := utils.MakeAuthCookie(userWithSession.SessionID, h.backendURI)
	http.SetCookie(w, cookie)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body: &dto.JSONUser{
			ID:       userWithSession.ID,
			UserType: userWithSession.UserType,
		},
	})
}
