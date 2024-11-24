package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"

	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
)

type CVsHandler struct {
	logger               *logrus.Entry
	cvsUsecase           cvs.ICVsUsecase
	sessionApplicantRepo session.ISessionRepository
	fileLoadingUsecase   fileloading.IFileLoadingUsecase
}

func NewCVsHandler(layers *internal.App) *CVsHandler {
	fn := "CVsHandler.NewCVsHandler"
	if layers == nil {
		logrus.Fatalf("%s: layers is nil", fn)
		return nil
	}
	logger := layers.Logger
	if layers.Usecases == nil || layers.Usecases.CVUsecase == nil {
		logrus.Fatalf("%s: layers is nil", fn)
		return nil
	}
	logger.Debug("CVsHandler created")
	return &CVsHandler{
		logger:               &logrus.Entry{Logger: logger},
		cvsUsecase:           layers.Usecases.CVUsecase,
		sessionApplicantRepo: layers.Repositories.SessionApplicantRepository,
		fileLoadingUsecase:   layers.Usecases.FileLoadingUsecase,
	}
}

func (h *CVsHandler) CVsRESTHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Logger.Debugf("CVsHandler.CVsRESTHandler got request: %s", r.URL.Path)
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
			Error:      dto.MsgMethodNotAllowed,
		})
		r.Body.Close()
	}
}

func (h *CVsHandler) CreateCVHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.CreateCVHandler"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	r.ParseMultipartForm(25 << 20) // 25Mb
	newCV := &dto.JSONCv{}
	newCV.PositionRu = r.FormValue("positionRu")
	newCV.PositionEn = r.FormValue("positionEn")
	newCV.Description = r.FormValue("description")
	newCV.JobSearchStatusName = r.FormValue("jobSearchStatusName")
	newCV.WorkingExperience = r.FormValue("workingExperience")
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
		newCV.Avatar = fileAddress
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
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
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

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
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, ID)
	cvID := uint64(ID)
	r.ParseMultipartForm(25 << 20) // 25Mb
	newCV := &dto.JSONCv{}
	newCV.PositionRu = r.FormValue("positionRu")
	newCV.PositionEn = r.FormValue("positionEn")
	newCV.Description = r.FormValue("description")
	newCV.JobSearchStatusName = r.FormValue("jobSearchStatusName")
	newCV.WorkingExperience = r.FormValue("workingExperience")
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
		newCV.Avatar = fileAddress
	}

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnauthorized,
		})
		return
	}

	updatedCV, err := h.cvsUsecase.UpdateCV(cvID, currentUser, newCV)
	if err != nil {
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
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, ID)

	cvID := uint64(ID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	err = h.cvsUsecase.DeleteCV(cvID, currentUser)
 	if err != nil {
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
