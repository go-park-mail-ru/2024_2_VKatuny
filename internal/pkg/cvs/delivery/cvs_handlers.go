package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

type CVsHandler struct {
	logger     *logrus.Logger
	cvsUsecase cvs.ICVsUsecase
}

func NewCVsHandler(logger *logrus.Logger, usecases *internal.Usecases) *CVsHandler {
	logger.Debug("CVsHandler created")
	return &CVsHandler{
		logger:     logger,
		cvsUsecase: usecases.CVUsecase,
	}
}

func (h *CVsHandler) CVsRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("CVsHandler.CVsRESTHandler got request: %s", r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		// TODO: навесить миддлвару проверки авторизации
		h.CreateCVHandler(w, r)
	case http.MethodGet:
		h.GetCVsHandler(w, r)
	case http.MethodPut:
		// TODO: навесить миддлвару проверки авторизации
		h.UpdateCVHandler(w, r)
	case http.MethodDelete:
		// TODO: навесить миддлвару проверки авторизации
		h.DeleteCVHandler(w, r)
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

	wroteCV, err := h.cvsUsecase.CreateCV(newCV)
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

	session, err := r.Cookie(dto.SessionIDName)
	if err == http.ErrNoCookie || session.Value == "" {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      http.ErrNoCookie.Error(),
		})
		return
	}

	updatedCV, err := h.cvsUsecase.UpdateCV(cvID, session.Value, newCV)
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

	session, err := r.Cookie(dto.SessionIDName)
	if err == http.ErrNoCookie || session.Value == "" {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      http.ErrNoCookie.Error(),
		})
		return
	}

	err = h.cvsUsecase.DeleteCV(cvID, session.Value)
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
