package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

// CreateApplicantHandler creates applicant in db
// CreateApplicant godoc
// @Summary     Creates a new user as a applicant
// @Description -
// @Tags        Applicant
// @Accept      json
// @Produce     json
// @Param       example body     dto.JSONApplicantRegistrationForm true "Example"
// @Success     200 {object} dto.JSONResponse{body=dto.JSONUser}
// @Failure     400 {object} dto.JSONResponse
// @Router      /api/v1/applicant/registration [post]
func (h *ApplicantHandlers) ApplicantRegistration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantRegistration"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

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
	applicantRegistrationForm.FirstName =
		utils.SanitizeString(applicantRegistrationForm.FirstName)
	applicantRegistrationForm.LastName =
		utils.SanitizeString(applicantRegistrationForm.LastName)
	applicantRegistrationForm.Email =
		utils.SanitizeString(applicantRegistrationForm.Email)
	applicantRegistrationForm.BirthDate =
		utils.SanitizeString(applicantRegistrationForm.BirthDate)
	applicantRegistrationForm.Password =
		utils.SanitizeString(applicantRegistrationForm.Password)
	h.logger.Debugf("%s: json decoded: %v", fn, applicantRegistrationForm)

	// TODO: add json validation with govalidator
	applicant, err := h.applicantUsecase.Create(r.Context(), applicantRegistrationForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgUnableToCreateUser,
		})
		return
	}
	h.logger.Debugf("%s: user created successfully: %v", fn, applicant)

	requestID := r.Context().Value(dto.RequestIDContextKey).(string)
	grpc_request := &grpc_auth.AuthRequest{
		RequestID: requestID,
		UserType:  dto.UserTypeApplicant,
		Email:     applicantRegistrationForm.Email,
		Password:  applicantRegistrationForm.Password,
	}
	grpc_response, err := h.authGRPC.AuthUser(r.Context(), grpc_request)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(), // TODO: standardize errors
		})
		return
	}

	user := &dto.JSONUser{
		ID:       grpc_response.UserData.ID,
		UserType: grpc_response.UserData.UserType,
	}
	user.UserType = utils.SanitizeString(user.UserType)

	h.logger.Debugf("%s: user logged in: %v", fn, user)

	cookie := utils.MakeAuthCookie(grpc_response.Session.ID, h.backendURI)
	http.SetCookie(w, cookie)
	h.logger.Debugf("%s: cookie set: %s", fn, cookie.Value)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       user,
	})
}
