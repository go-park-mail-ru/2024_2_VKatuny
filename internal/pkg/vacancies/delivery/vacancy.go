package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type VacanciesHandlers struct {
	logger               *logrus.Entry
	vacanciesUsecase     vacancies.IVacanciesUsecase
	fileLoadingUsecase   fileloading.IFileLoadingUsecase
	CompressGRPC         compressmicroservice.CompressServiceClient
}

func NewVacanciesHandlers(layers *internal.App) *VacanciesHandlers {
	logger := layers.Logger
	logger.Debug("VacanciesHandlers created")

	return &VacanciesHandlers{
		logger:               &logrus.Entry{Logger: logger},
		vacanciesUsecase:     layers.Usecases.VacanciesUsecase,
		fileLoadingUsecase:   layers.Usecases.FileLoadingUsecase,
		CompressGRPC:         layers.Microservices.Compress,
	}
}

// @Summary CreateVacancy
// @Description Create new vacancy
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   newVacancy body     dto.JSONVacancy true "New vacancy"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Router /api/v1/vacancy [post]
func (h *VacanciesHandlers) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	fn := "VacanciesHandlers.CreateVacancy"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	r.ParseMultipartForm(25 << 20) // 25Mb
	newVacancy := &dto.JSONVacancy{}
	newVacancy.Position = r.FormValue("position")
	newVacancy.Location = r.FormValue("location")
	newVacancy.Description = r.FormValue("description")
	newVacancy.WorkType = r.FormValue("workType")
	newVacancy.CompanyName = r.FormValue("companyName")
	newVacancy.PositionCategoryName = r.FormValue("group")
	temp, err := strconv.Atoi(r.FormValue("salary"))
	utils.EscapeHTMLStruct(newVacancy)
	if err != nil {
		h.logger.Errorf("bad input of salary: %s", err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
		return
	}
	newVacancy.Salary = int32(temp)
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("company_avatar")
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
		newVacancy.Avatar = fileAddress
		newVacancy.CompressedAvatar = compressedFileAddress
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	h.logger.Debug(newVacancy)
	vacancy, err := h.vacanciesUsecase.CreateVacancy(newVacancy, currentUser)
	if err != nil {
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       vacancy,
	})
}

// @Summary GetVacancy
// @Description Get vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id} [get]
func (h *VacanciesHandlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.GetVacancy"

	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	vacancy, err := h.vacanciesUsecase.GetVacancy(vacancyID)
	if err != nil {
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got vacancy: %v", fn, vacancy)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       vacancy,
	})
}

// @Summary UpdateVacancy
// @Description Update vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Param   position formData string true "Vacancy position"
// @Param   location formData string true "Vacancy location"
// @Param   description formData string true "Vacancy description"
// @Param   workType formData string true "Vacancy workType"
// @Param   companyName formData string true "Company name"
// @Param   group formData string true "Vacancy group"
// @Param   salary formData int32 true "Vacancy salary"
// @Param   company_avatar formData file true "Company avatar"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id} [put]
func (h *VacanciesHandlers) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}
	fn := "VacanciesHandlers.UpdateVacancy"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	r.ParseMultipartForm(25 << 20) // 25Mb
	updatedVacancy := &dto.JSONVacancy{}
	updatedVacancy.Position = r.FormValue("position")
	updatedVacancy.Location = r.FormValue("location")
	updatedVacancy.Description = r.FormValue("description")
	updatedVacancy.WorkType = r.FormValue("workType")
	updatedVacancy.CompanyName = r.FormValue("companyName")
	updatedVacancy.PositionCategoryName = r.FormValue("group")
	temp, err := strconv.Atoi(r.FormValue("salary"))
	if err != nil {
		h.logger.Errorf("bad input of salary: %s", err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
		return
	}
	updatedVacancy.Salary = int32(temp)
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("company_avatar")
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
		updatedVacancy.Avatar = fileAddress
		updatedVacancy.CompressedAvatar = compressedFileAddress
	}
	
	wroteVacancy, err := h.vacanciesUsecase.UpdateVacancy(vacancyID, updatedVacancy, currentUser)
	if err != nil {
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       wroteVacancy,
	})
}

// @Summary DeleteVacancy
// @Description Delete vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id} [delete]
func (h *VacanciesHandlers) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.DeleteVacancy"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusNotFound, dto.JSONResponse{
			HTTPStatus: http.StatusNotFound,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	err = h.vacanciesUsecase.DeleteVacancy(vacancyID, currentUser)
	if err != nil {
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

// @Summary Subscribe on vacancy
// @Description Subscribe on vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id}/subscription [post]
func (h *VacanciesHandlers) SubscribeVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.SubscribeVacancy"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	err = h.vacanciesUsecase.SubscribeOnVacancy(vacancyID, currentUser)
	if err != nil {
		h.logger.Errorf("while subscribing on vacancy got: %s", err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("user_ID: %d subscribed on vacancy_ID %d", currentUser.ID, vacancyID)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

// @Summary Unsubscribe from vacancy
// @Description Unsubscribe from vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id}/subscription [delete]
func (h *VacanciesHandlers) UnsubscribeVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.UnsubscribeVacancy"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	err = h.vacanciesUsecase.UnsubscribeFromVacancy(vacancyID, currentUser)
	if err != nil {
		h.logger.Errorf("while unsubscribing from vacancy got: %s", err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("user_ID: %d unsubscribed from vacancy_ID %d", currentUser.ID, vacancyID)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

// @Summary Get subscription status on vacancy
// @Description Get subscription status on vacancy by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONVacancySubscriptionStatus
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id}/subscription [get]
func (h *VacanciesHandlers) GetVacancySubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.GetVacancySubscription"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	vacancySubscriptionInfo, err := h.vacanciesUsecase.GetSubscriptionInfo(vacancyID, currentUser.ID)
	if err != nil {
		h.logger.Errorf("while getting subscription status got: %s", err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       vacancySubscriptionInfo,
	})
}

// @Summary Get vacancy subscribers
// @Description Get vacancy subscribers by ID
// @Tags Vacancy
// @Accept  json
// @Produce  json
// @Param   id path string true "Vacancy ID"
// @Success 200 {object} dto.JSONVacancySubscribers
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/vacancy/{id}/subscribers [get]
func (h *VacanciesHandlers) GetVacancySubscribers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "VacanciesHandlers.GetVacancySubscribers"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	vacancyID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug: %d", fn, vacancyID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	subscribers, err := h.vacanciesUsecase.GetVacancySubscribers(vacancyID, currentUser)
	if err != nil {
		h.logger.Errorf("while getting subscribers got: %s", err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       subscribers,
	})
}
