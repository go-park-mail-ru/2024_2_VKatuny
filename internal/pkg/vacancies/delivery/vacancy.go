package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type VacanciesHandlers struct {
	logger               *logrus.Entry
	vacanciesUsecase     vacancies.IVacanciesUsecase
	sessionEmployerRepo  session.ISessionRepository
	sessionApplicantRepo session.ISessionRepository
	fileLoadingUsecase   fileloading.IFileLoadingUsecase
}

func NewVacanciesHandlers(layers *internal.App) *VacanciesHandlers {
	logger := layers.Logger
	logger.Debug("VacanciesHandlers created")

	return &VacanciesHandlers{
		logger:               &logrus.Entry{Logger: logger},
		vacanciesUsecase:     layers.Usecases.VacanciesUsecase,
		sessionEmployerRepo:  layers.Repositories.SessionEmployerRepository,
		sessionApplicantRepo: layers.Repositories.SessionApplicantRepository,
		fileLoadingUsecase:   layers.Usecases.FileLoadingUsecase,
	}
}

func (h *VacanciesHandlers) VacanciesRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Logger.Debugf("VacanciesHandlers.VacanciesRESTHandler got request: %s", r.URL.Path)
	repository := &internal.Repositories{SessionEmployerRepository: h.sessionEmployerRepo}
	switch r.Method {
	case http.MethodPost:
		handler := middleware.RequireAuthorization(h.createVacancyHandler, repository, dto.UserTypeEmployer)
		handler(w, r)
	case http.MethodGet:
		h.getVacancyHandler(w, r)
	case http.MethodPut:
		handler := middleware.RequireAuthorization(h.updateVacancyHandler, repository, dto.UserTypeEmployer)
		handler(w, r)
	case http.MethodDelete:
		handler := middleware.RequireAuthorization(h.deleteVacancyHandler, repository, dto.UserTypeEmployer)
		handler(w, r)
	default:
		middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
			HTTPStatus: http.StatusMethodNotAllowed,
			Error:      dto.MsgMethodNotAllowed,
		})
	}
}

func (h *VacanciesHandlers) VacanciesSubscribeRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Logger.Debugf("VacanciesHandlers.VacanciesSubscribeRESTHandler got request: %s", r.URL.Path)
	repository := &internal.Repositories{SessionApplicantRepository: h.sessionApplicantRepo}
	switch r.Method {
	case http.MethodPost:
		handler := middleware.RequireAuthorization(h.subscribeVacancyHandler, repository, dto.UserTypeApplicant)
		handler(w, r)
	case http.MethodGet:
		handler := middleware.RequireAuthorization(h.getVacancySubscriptionHandler, repository, dto.UserTypeApplicant)
		handler(w, r)
	case http.MethodDelete:
		handler := middleware.RequireAuthorization(h.unsubscribeVacancyHandler, repository, dto.UserTypeApplicant)
		handler(w, r)
	default:
		middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
			HTTPStatus: http.StatusMethodNotAllowed,
			Error:      dto.MsgMethodNotAllowed,
		})
	}
}

func (h *VacanciesHandlers) GetVacancySubscribersHandler(w http.ResponseWriter, r *http.Request) {
	repository := &internal.Repositories{SessionEmployerRepository: h.sessionEmployerRepo}
	handler := middleware.RequireAuthorization(h.getVacancySubscribersHandler, repository, dto.UserTypeEmployer)
	handler(w, r)
}

func (h *VacanciesHandlers) createVacancyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(25 << 20) // 25Mb
	newVacancy := &dto.JSONVacancy{}
	newVacancy.Position = r.FormValue("position")
	newVacancy.Location = r.FormValue("location")
	newVacancy.Description = r.FormValue("description")
	newVacancy.WorkType = r.FormValue("workType")
	newVacancy.CompanyName = r.FormValue("companyName")
	temp, err := strconv.Atoi(r.FormValue("salary"))
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

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
		fileAddress, err := h.fileLoadingUsecase.WriteImage(file, header)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		newVacancy.Avatar = fileAddress
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	h.logger.Debug(newVacancy)
	wroteVacancy, err := h.vacanciesUsecase.CreateVacancy(newVacancy, currentUser)
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

func (h *VacanciesHandlers) getVacancyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	vacancy, err := h.vacanciesUsecase.GetVacancy(vacancyID)
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

func (h *VacanciesHandlers) updateVacancyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	r.ParseMultipartForm(25 << 20) // 25Mb
	updatedVacancy := &dto.JSONVacancy{}
	updatedVacancy.Position = r.FormValue("position")
	updatedVacancy.Location = r.FormValue("location")
	updatedVacancy.Description = r.FormValue("description")
	updatedVacancy.WorkType = r.FormValue("workType")
	updatedVacancy.CompanyName = r.FormValue("companyName")
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
		fileAddress, err := h.fileLoadingUsecase.WriteImage(file, header)
		if err != nil {
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		updatedVacancy.Avatar = fileAddress
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
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

func (h *VacanciesHandlers) deleteVacancyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
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

func (h *VacanciesHandlers) subscribeVacancyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/subscription/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
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

func (h *VacanciesHandlers) unsubscribeVacancyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/subscription/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
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

func (h *VacanciesHandlers) getVacancySubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/subscription/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
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

func (h *VacanciesHandlers) getVacancySubscribersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/subscribers/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error(dto.MsgUnableToGetUserFromContext)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	vacancyID := uint64(slug)

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
