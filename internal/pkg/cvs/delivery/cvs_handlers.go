package delivery

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
)

type CVsHandler struct {
	logger             *logrus.Entry
	applicantUsecase   applicant.IApplicantUsecase
	cvsUsecase         cvs.ICVsUsecase
	fileLoadingUsecase fileloading.IFileLoadingUsecase
	CompressGRPC       compressmicroservice.CompressServiceClient
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
		logger:             &logrus.Entry{Logger: logger},
		applicantUsecase:   layers.Usecases.ApplicantUsecase,
		cvsUsecase:         layers.Usecases.CVUsecase,
		fileLoadingUsecase: layers.Usecases.FileLoadingUsecase,
		CompressGRPC:       layers.Microservices.Compress,
	}
}

// @Summary Create CV
// @Description Create new CV
// @Tags CV
// @Accept  multipart/form-data
// @Produce  json
// @Param   profile_avatar     formData  file     true  "Profile avatar"
// @Param   positionRu         formData  string   true  "Position in Russian"
// @Param   positionEn         formData  string   true  "Position in English"
// @Param   description        formData  string   true  "Description"
// @Param   jobSearchStatus formData  string   true  "Job search status name"
// @Param   workingExperience   formData  string   true  "Working experience"
// @Param   group              formData  string   true  "Group"
// @Success 200 {object} dto.JSONResponse{body=dto.JSONCv}
// @Failure 400 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/cv [POST]
func (h *CVsHandler) CreateCV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.CreateCVHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	r.ParseMultipartForm(25 << 20) // 25Mb
	newCV := &dto.JSONCv{}
	newCV.PositionRu = r.FormValue("positionRu")
	newCV.PositionEn = r.FormValue("positionEn")
	newCV.Description = r.FormValue("description")
	newCV.JobSearchStatusName = r.FormValue("jobSearchStatus")
	newCV.WorkingExperience = r.FormValue("workingExperience")
	newCV.PositionCategoryName = r.FormValue("group")
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
		newCV.Avatar = fileAddress
		newCV.CompressedAvatar = compressedFileAddress
	}
	utils.EscapeHTMLStruct(newCV)
	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	wroteCV, err := h.cvsUsecase.CreateCV(r.Context(), newCV, currentUser)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: success, got created cv: %v", fn, wroteCV)

	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       wroteCV,
	})
}

// GetCV godoc
// @Summary     Gets CV by id
// @Description Gets CV by id
// @Tags        CV
// @Accept      json
// @Produce     json
// @Param       id   path      uint64  true  "CV ID"
// @Success     200  {object}  dto.JSONGetApplicantCV
// @Failure     400  {object}  dto.JSONResponse
// @Failure     405  {object}  dto.JSONResponse
// @Failure     500  {object}  dto.JSONResponse
// @Router      /api/v1/cv/{id} [get]
func (h *CVsHandler) GetCV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.GetCVsHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s; entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	cvID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got slug cvID: %d", fn, cvID)

	CV, err := h.cvsUsecase.GetCV(r.Context(), cvID)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	CV.CompressedAvatar  = h.fileLoadingUsecase.FindCompressedFile(CV.Avatar)
	h.logger.Debugf("%s: success, got cv: %v", fn, CV)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       CV,
	})
}

// @Tags CV
// @Summary Update CV
// @Security ApiKeyAuth
// @Description Update CV
// @Accept  json
// @Produce  json
// @Param   id   path     uint64     true  "CV ID"
// @Param   positionRu     formData   string     true  "Position in Russian"
// @Param   positionEn     formData   string     true  "Position in English"
// @Param   description     formData   string     true  "Description"
// @Param   jobSearchStatus     formData   string     true  "Job search status name"
// @Param   workingExperience     formData   string     true  "Working experience"
// @Param   group     formData   string     true  "Group"
// @Param   profile_avatar     formData   file     true  "Profile avatar"
// @Success 200 {object} dto.JSONResponse{body=[]dto.JSONCv}
// @Failure 400 {object} dto.JSONResponse
// @Failure 401 {object} dto.JSONResponse
// @Failure 403 {object} dto.JSONResponse
// @Failure 404 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router      /api/v1/cv/{id} [put]
func (h *CVsHandler) UpdateCV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.UpdateCVHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	cvID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, cvID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnauthorized,
		})
		return
	}

	r.ParseMultipartForm(25 << 20) // 25Mb
	newCV := &dto.JSONCv{}
	newCV.PositionRu = r.FormValue("positionRu")
	newCV.PositionEn = r.FormValue("positionEn")
	newCV.Description = r.FormValue("description")
	newCV.JobSearchStatusName = r.FormValue("jobSearchStatus")
	newCV.WorkingExperience = r.FormValue("workingExperience")
	newCV.PositionCategoryName = r.FormValue("group")
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
		newCV.Avatar = fileAddress
		newCV.CompressedAvatar = compressedFileAddress
	}
	utils.EscapeHTMLStruct(newCV)
	updatedCV, err := h.cvsUsecase.UpdateCV(r.Context(), cvID, currentUser, newCV)
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

// @Tags CV
// @Summary Delete cv
// @Description Delete cv
// @Accept json
// @Produce json
// @Param id path uint64 true "id of cv"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 401 {object} dto.JSONResponse
// @Failure 404 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/cv/{id} [delete]
func (h *CVsHandler) DeleteCV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.DeleteCVHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	cvID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, cvID)

	currentUser, ok := r.Context().Value(dto.UserContextKey).(*dto.UserFromSession)
	if !ok {
		h.logger.Error("unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization")
		middleware.UniversalMarshal(w, http.StatusUnauthorized, dto.JSONResponse{
			HTTPStatus: http.StatusUnauthorized,
			Error:      dto.MsgUnableToGetUserFromContext,
		})
		return
	}

	err = h.cvsUsecase.DeleteCV(r.Context(), cvID, currentUser)
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

// @Tags CV
// @Summary Get cv
// @Description Get cv
// @Accept slug
// @Produce json
// @Param id path uint64 true "id of cv"
// @Success 200 {object} dto.JSONResponse
// @Failure 400 {object} dto.JSONResponse
// @Failure 401 {object} dto.JSONResponse
// @Failure 404 {object} dto.JSONResponse
// @Failure 405 {object} dto.JSONResponse
// @Failure 500 {object} dto.JSONResponse
// @Router /api/v1/cv-to-pdf/{id:[0-9]+} [get]
func (h *CVsHandler) CVtoPDF(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "CVsHandler.CVtoPDF"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)
	h.logger.Debugf("%s: entering", fn)

	vars := mux.Vars(r)
	slug := vars["id"]
	cvID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: got slug cvID: %d", fn, cvID)

	var name *dto.CVPDFFile
	CV, err := h.cvsUsecase.GetCV(r.Context(), cvID)
	if err == nil {
		applicant, err := h.applicantUsecase.GetApplicantProfile(context.Background(), CV.ApplicantID)
		if err == nil {
			name, err = h.fileLoadingUsecase.CVtoPDF(CV, applicant)
		}
	}
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	name.FileName="/"+name.FileName
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       name,
	})
}
