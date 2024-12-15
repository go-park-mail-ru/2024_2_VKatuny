package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

type EmployerHandlers struct {
	logger             *logrus.Entry
	backendAddress     string
	employerUsecase    employer.IEmployerUsecase
	vacanciesUsecase   vacancies.IVacanciesUsecase
	fileLoadingUsecase fileloading.IFileLoadingUsecase
	authGRPC           auth_grpc.AuthorizationClient
	compressGRPC       compressmicroservice.CompressServiceClient
}

func NewEmployerHandlers(app *internal.App) *EmployerHandlers {
	return &EmployerHandlers{
		logger:             logrus.NewEntry(app.Logger),
		backendAddress:     app.BackendAddress,
		employerUsecase:    app.Usecases.EmployerUsecase,
		vacanciesUsecase:   app.Usecases.VacanciesUsecase,
		fileLoadingUsecase: app.Usecases.FileLoadingUsecase,
		authGRPC:           app.Microservices.Auth,
		compressGRPC:       app.Microservices.Compress,
	}
}

// @Summary Get employer profile
// @Description Get employer profile by ID
// @Tags Employer
// @Accept json
// @Produce json
// @Param id path string true "Employer ID"
// @Success 200 {object} dto.JSONResponse{body=dto.JSONGetEmployerProfile}
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/employer/{id}/profile [get]
func (h *EmployerHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerProfileHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)

	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	// dto - JSONGetEmployerProfile
	profile, err := h.employerUsecase.GetEmployerProfile(r.Context(), employerID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	utils.EscapeHTMLStruct(profile)

	h.logger.Debugf("function %s: success, got profile %v", fn, profile)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       profile,
	})
}

// @Summary Update employer profile
// @Description Update employer profile by ID
// @Tags Employer
// @Accept multipart/form-data
// @Produce json
// @Param id path uint64 true "Employer ID"
// @Param firstName formData string true "First Name"
// @Param lastName formData string true "Last Name"
// @Param city formData string true "City"
// @Param contacts formData string true "Contacts"
// @Param profile_avatar formData file false "Profile Avatar"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/employer/{id}/profile [put]
func (h *EmployerHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.UpdateEmployerProfileHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)
	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	newProfileData := &dto.JSONUpdateEmployerProfile{}
	newProfileData.FirstName = r.FormValue("firstName")
	newProfileData.LastName = r.FormValue("lastName")
	newProfileData.City = r.FormValue("city")
	newProfileData.Contacts = r.FormValue("contacts")
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("profile_avatar")
	if err == nil {
		defer file.Close()
		fileAddress, compressedFileAddress, err := h.fileLoadingUsecase.WriteImage(file, header)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		newProfileData.Avatar = fileAddress
		newProfileData.CompressedAvatar = compressedFileAddress
	}
	utils.EscapeHTMLStruct(newProfileData)
	h.logger.Debugf("function %s: new profile data MultiPart parsed: %v", fn, newProfileData)

	err = h.employerUsecase.UpdateEmployerProfile(r.Context(), employerID, newProfileData)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: success", fn)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

// @Summary Get employer vacancies
// @Description Get vacancies by employer ID
// @Tags Employer
// @Accept json
// @Produce json
// @Param id path string true "Employer ID"
// @Success 200 {object} dto.JSONResponse{body=[]models.Vacancy}
// @Failure 405 {object} dto.JSONResponse
// @Router /api/v1/employer/{id}/vacancies [get]
func (h *EmployerHandlers) GetEmployerVacancies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerVacanciesHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)
	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	vacancies, err := h.vacanciesUsecase.GetVacanciesByEmployerID(employerID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	for _, vacancy := range vacancies {
		utils.EscapeHTMLStruct(vacancy)
	} 

	h.logger.Debugf("function %s: success, got vacancies: %d", fn, len(vacancies))
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       vacancies,
	})
}
