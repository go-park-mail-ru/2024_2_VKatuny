// Package delivery is a handlers layer of employer
package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

// CreateEmployerHandler creates employers in db
// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Employer
// @Accept      json
// @Produce     json
// @Param       example body     dto.JSONEmployerRegistrationForm true "Example"
// @Success     200 {object} dto.JSONResponse{body=dto.JSONUser}
// @Failure     400 {object} dto.JSONResponse
// @Failure     405 {object} dto.JSONResponse
// @Failure     500 {object} dto.JSONResponse
// @Router      /api/v1/employer/registration [post]
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

	utils.EscapeHTMLStruct(employerRegistrationForm)
	h.logger.Debugf("%s: json decoded successfully: %v", fn, employerRegistrationForm)

	// TODO: implement usecase for validate registration data

	_, err = h.employerUsecase.Create(r.Context(), employerRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: employer created successfully", fn)

	requestID := r.Context().Value(dto.RequestIDContextKey).(string)
	grpc_request := &grpc_auth.AuthRequest{
		RequestID: requestID,
		UserType:  dto.UserTypeEmployer,
		Email:     employerRegistrationForm.Email,
		Password:  employerRegistrationForm.Password,
	}
	grpc_response, err := h.authGRPC.AuthUser(r.Context(), grpc_request)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	user := &dto.JSONUser{
		ID:       grpc_response.UserData.ID,
		UserType: grpc_response.UserData.UserType,
	}
	user.UserType = utils.EscapeHTMLString(user.UserType)
	h.logger.Debugf("%s: employer logged in successfully: %v", fn, user)

	cookie := utils.MakeAuthCookie(grpc_response.Session.ID, h.backendAddress)
	h.logger.Debugf("%s: cookie created %s", fn, cookie.Value)
	http.SetCookie(w, cookie)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       user,
	})
}
