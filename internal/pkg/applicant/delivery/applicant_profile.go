package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	auth_grpc "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
)

type ApplicantHandlers struct {
	logger             *logrus.Entry
	backendURI         string
	applicantUsecase   applicant.IApplicantUsecase
	portfolioUsecase   portfolio.IPortfolioUsecase
	cvUsecase          cvs.ICVsUsecase
	fileLoadingUsecase fileloading.IFileLoadingUsecase
	authGRPC           auth_grpc.AuthorizationClient
	CompressGRPC       compressmicroservice.CompressServiceClient
}

func NewApplicantProfileHandlers(app *internal.App) *ApplicantHandlers {
	return &ApplicantHandlers{
		logger:             logrus.NewEntry(app.Logger),
		backendURI:         app.BackendAddress,
		applicantUsecase:   app.Usecases.ApplicantUsecase,
		portfolioUsecase:   app.Usecases.PortfolioUsecase,
		cvUsecase:          app.Usecases.CVUsecase,
		fileLoadingUsecase: app.Usecases.FileLoadingUsecase,
		authGRPC:           app.Microservices.Auth,
		CompressGRPC:       app.Microservices.Compress,
	}
}

// GetProfile godoc
// @Summary Get applicant profile
// @Description Get applicant profile by ID
// @Tags Applicant
// @Accept json
// @Produce json
// @Param id path string true "Applicant ID"
// @Success 200 {object} dto.JSONGetApplicantProfile
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/applicant/{id}/profile [get]
func (h *ApplicantHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantProfileHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)

	slug := vars["id"]
	applicantID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	// dto - JSONGetApplicantProfile
	applicantProfile, err := h.applicantUsecase.GetApplicantProfile(r.Context(), applicantID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got profile %v", fn, applicantProfile)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       applicantProfile,
	})
}

// @Tags Applicant
// @Summary Update applicant profile
// @Description Update applicant profile
// @Accept json
// @Produce json
// @Param id path uint64 true "ID of applicant"
// @Param input body dto.JSONUpdateApplicantProfile true "Profile to update"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/applicant/{id}/profile [put]
func (h *ApplicantHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.UpdateApplicantProfileHandler"

	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)

	slug := vars["id"]
	applicantID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	newProfileData := &dto.JSONUpdateApplicantProfile{}
	newProfileData.FirstName = r.FormValue("firstName")
	newProfileData.LastName = r.FormValue("lastName")
	newProfileData.City = r.FormValue("city")
	newProfileData.BirthDate = r.FormValue("birthDate")
	newProfileData.Contacts = r.FormValue("contacts")
	newProfileData.Education = r.FormValue("education")
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("profile_avatar")
	if err == nil {
		defer file.Close()
		fileAddress, err := h.fileLoadingUsecase.WriteImage(file, header)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		newProfileData.Avatar = fileAddress
		var buff []byte
		file.Read(buff)
		h.logger.Debugf("Start compression")
		_, err = h.CompressGRPC.CompressAndSaveFile(
			r.Context(),
			&compressmicroservice.CompressAndSaveFileInput{
				FileName: fileAddress + header.Filename,
				FileType: header.Header["Content-Type"][0],
				File:     buff,
			},
		)
	}
	err = h.applicantUsecase.UpdateApplicantProfile(r.Context(), applicantID, newProfileData)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: successfully updated profile", fn)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

// GetPortfolios godoc
// @Summary Get applicant portfolios
// @Description Get portfolios of an applicant by ID
// @Tags Applicant
// @Accept json
// @Produce json
// @Param id path string true "Applicant ID"
// @Success 200 {object} dto.JSONResponse "portfolios"
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/applicant/{id}/portfolio [get]
func (h *ApplicantHandlers) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantPortfoliosHandler"

	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)

	slug := vars["id"]
	applicantID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	portfolios, err := h.portfolioUsecase.GetApplicantPortfolios(r.Context(), applicantID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got portfolios: %d", fn, len(portfolios))
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       portfolios,
	})
}

// GetCVs godoc
// @Summary Get applicant CVs
// @Description Get CVs of an applicant by ID
// @Tags Applicant
// @Accept json
// @Produce json
// @Param id path string true "Applicant ID"
// @Success 200 {object} dto.JSONResponse "CVs"
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/applicant/{id}/cv [get]
func (h *ApplicantHandlers) GetCVs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantCVsHandler"

	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)

	slug := vars["id"]
	applicantID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	// *dto.JSONGetApplicantCV
	CVs, err := h.cvUsecase.GetApplicantCVs(applicantID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got CVs: %d", fn, len(CVs))
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       CVs,
	})
}
