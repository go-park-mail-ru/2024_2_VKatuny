package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/sirupsen/logrus"
)

type CVsHandler struct {
	logger               *logrus.Logger
	cvsUsecase           cvs.ICVsUsecase
	sessionApplicantRepo session.ISessionRepository
}

func NewCVsHandler(layers *internal.App) *CVsHandler {
	logger := layers.Logger
	logger.Debug("CVsHandler created")
	return &CVsHandler{
		logger:     logger,
		cvsUsecase: layers.Usecases.CVUsecase,
		sessionApplicantRepo: layers.Repositories.SessionApplicantRepository,
	}
}

func (h *CVsHandler) CVsRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("CVsHandler.CVsRESTHandler got request: %s", r.URL.Path)
	repositories := &internal.Repositories{SessionApplicantRepository: h.sessionApplicantRepo}
	switch r.Method {
	case http.MethodPost:
		handler := middleware.RequireAuthorization(h.CreateCVHandler, repositories, dto.UserTypeApplicant)
		handler(w, r)
	case http.MethodGet:
		h.GetCVsHandler(w, r)
	case http.MethodPut:
		handler := middleware.RequireAuthorization(h.UpdateCVHandler, repositories, dto.UserTypeApplicant)
		handler(w, r)
	case http.MethodDelete:
		handler := middleware.RequireAuthorization(h.DeleteCVHandler, repositories, dto.UserTypeApplicant)
		handler(w, r)
	default:
		middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
			HTTPStatus: http.StatusMethodNotAllowed,
			Error:      http.StatusText(http.StatusMethodNotAllowed),
		})
		r.Body.Close()
	}
}

func (h *CVsHandler) CreateCVHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.CreateCVHandler"

	decoder := json.NewDecoder(r.Body)
	newCV := new(dto.JSONCv)

	err := decoder.Decode(newCV)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      commonerrors.ErrInvalidJSON.Error(),
		})
		return
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.SessionUser)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      "unable to get user from context", // TODO: make error without hardcode
		})
		return
	}

	wroteCV, err := h.cvsUsecase.CreateCV(newCV, currentUser)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: success, got created cv: %v", fn, wroteCV)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       wroteCV,
	})
}

func (h *CVsHandler) GetCVsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.GetCVsHandler"

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, ID)

	cvID := uint64(ID)

	CV, err := h.cvsUsecase.GetCV(cvID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got cv: %v", fn, CV)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       CV,
	})
}

func (h *CVsHandler) UpdateCVHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fn := "CVsHandler.UpdateCVHandler"
	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, ID)
	cvID := uint64(ID)
	decoder := json.NewDecoder(r.Body)
	newCV := new(dto.JSONCv)
	err = decoder.Decode(newCV)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      commonerrors.ErrInvalidJSON.Error(),
		})
		return
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.SessionUser)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      "unable to get user from context", // TODO: make error without hardcode
		})
		return
	}

	updatedCV, err := h.cvsUsecase.UpdateCV(cvID, currentUser, newCV)
	if err == commonerrors.ErrUnauthorized || err == commonerrors.ErrSessionNotFound {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      err.Error(),
		})
		return
	} else if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got updated cv: %v", fn, updatedCV)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       updatedCV,
	})
}

func (h *CVsHandler) DeleteCVHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.DeleteCVHandler"
	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, ID)

	cvID := uint64(ID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.SessionUser)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      "unable to get user from context", // TODO: make error without hardcode
		})
		return
	}

	err = h.cvsUsecase.DeleteCV(cvID, currentUser)
	if err == commonerrors.ErrUnauthorized || err == commonerrors.ErrSessionNotFound {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      err.Error(),
		})
		return
	} else if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
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