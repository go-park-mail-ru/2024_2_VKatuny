// Package delivery is a handlers layer of employer
package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// CreateEmployerHandler creates employers in db
// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Success     200      {object}       dto.JSONResponse{statusCode=200,body=dto.JSONUserBody, error=""} "OK"
// @Failure     400      {object}       nil
// @Router      /registration/employer/ [post]
func (h *EmployerHandlers) Registration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerHandlers.CreateEmployerHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	employerRegistrationForm := new(dto.JSONEmployerRegistrationForm)
	err := json.NewDecoder(r.Body).Decode(employerRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
		return
	}
	h.logger.Debugf("%s: json decoded successfully: %v", fn, employerRegistrationForm)

	// TODO: implement usecase for validate registration data

	employer, err := h.employerUsecase.Create(r.Context(), employerRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: employer created successfully", fn)

	employerLogin := &dto.JSONLoginForm{
		UserType: dto.UserTypeEmployer,
		Email:    employerRegistrationForm.Email,
		Password: employerRegistrationForm.Password,
	}
	employerWithSession, err := h.sessionUsecase.Login(employerLogin)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: employer logged in successfully", fn)

	cookie := utils.MakeAuthCookie(employerWithSession.SessionID, h.backendAddress)
	h.logger.Debugf("%s: cookie created %s", fn, cookie.Value)
	http.SetCookie(w, cookie)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body: &dto.JSONUser{
			ID:       employer.ID,
			UserType: dto.UserTypeEmployer,
		},
	})
}
