package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/sirupsen/logrus"
)

type VacanciesHandlers struct {
	logger               *logrus.Logger
	vacanciesUsecase     vacancies.IVacanciesUsecase
	sessionEmployerRepo  session.ISessionRepository
	sessionApplicantRepo session.ISessionRepository
}

func NewVacanciesHandlers(layers *internal.App) *VacanciesHandlers {
	logger := layers.Logger
	logger.Debug("VacanciesHandlers created")

	return &VacanciesHandlers{
		logger:               logger,
		vacanciesUsecase:     layers.Usecases.VacanciesUsecase,
		sessionEmployerRepo:  layers.Repositories.SessionEmployerRepository,
		sessionApplicantRepo: layers.Repositories.SessionApplicantRepository,
	}
}

func (h *VacanciesHandlers) VacanciesRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("VacanciesHandlers.VacanciesRESTHandler got request: %s", r.URL.Path)
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
	h.logger.Debugf("VacanciesHandlers.VacanciesSubscribeRESTHandler got request: %s", r.URL.Path)
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
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	newVacancy := new(dto.JSONVacancy)

	err := decoder.Decode(newVacancy)
	if err != nil {
		h.logger.Errorf("unable to unmarshal JSON: %s", err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
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

	slug, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/vacancy/")
	if err != nil {
		h.logger.Errorf("while cutting slug got: %s", err)
		return
	}

	vacancyID := uint64(slug)

	decoder := json.NewDecoder(r.Body)
	updatedVacancy := new(dto.JSONVacancy)
	err = decoder.Decode(updatedVacancy)
	if err != nil {
		h.logger.Errorf("unable to unmarshal JSON: %s", err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      dto.MsgInvalidJSON,
		})
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
