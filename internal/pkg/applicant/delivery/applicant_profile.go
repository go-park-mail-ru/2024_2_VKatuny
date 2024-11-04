package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	"github.com/sirupsen/logrus"
)

type ApplicantProfileHandlers struct {
	logger           *logrus.Logger
	applicantUsecase applicantUsecase.IApplicantUsecase
	portfolioUsecase portfolioUsecase.IPortfolioUsecase
	cvUsecase        cvUsecase.ICVsUsecase 
}

func (h *ApplicantProfileHandlers) GetApplicantProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantProfileHandler"

	url := r.URL.Path[len("/api/v1/profile/applicant"):]
	slugID := strings.Split(url, "/")[0]
	ID, err := strconv.Atoi(slugID)
	if len(url) > 1 || err != nil {
		h.logger.Errorf("function %s: something bad with slug", fn)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "something bad with slug",
		})
		return
	}
	
	applicantID := uint64(ID)
	// dto - JSONGetApplicantProfile
	applicantProfile, err := h.applicantUsecase.GetApplicantProfile(applicantID)
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

func (h *ApplicantProfileHandlers) UpdateApplicantProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.UpdateApplicantProfileHandler"

	url := r.URL.Path[len("/api/v1/profile/applicant"):]
	slugID := strings.Split(url, "/")[0]
	ID, err := strconv.Atoi(slugID)
	if len(url) > 1 || err != nil {
		h.logger.Errorf("function %s: something bad with slug", fn)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "something bad with slug",
		})
		return
	}
	
	applicantID := uint64(ID)

	decoder := json.NewDecoder(r.Body)
	newProfileData := new(dto.JSONUpdateApplicantProfile)
	err = decoder.Decode(newProfileData)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "unable to unmarshal JSON",
		})
		return
	}

	err = h.applicantUsecase.UpdateApplicantProfile(applicantID, newProfileData)
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

func (h *ApplicantProfileHandlers) GetApplicantPortfoliosHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantPortfoliosHandler"

	url := r.URL.Path[len("/api/v1/portfolio/applicant/"):]
	slugID := strings.Split(url, "/")[0]
	ID, err := strconv.Atoi(slugID)
	if len(url) > 1 || err != nil {
		h.logger.Errorf("function %s: something bad with slug", fn)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "something bad with slug",
		})
		return
	}
	
	applicantID := uint64(ID)
	portfolios, err := h.portfolioUsecase.GetApplicantPortfolios(applicantID)
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

func (h *ApplicantProfileHandlers) GetApplicantCVsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantCVsHandler"

	url := r.URL.Path[len("/api/v1/cv/applicant/"):]
	slugID := strings.Split(url, "/")[0]
	ID, err := strconv.Atoi(slugID)
	if len(url) > 1 || err != nil {
		h.logger.Errorf("function %s: something bad with slug", fn)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "something bad with slug",
		})
		return
	}
	
	applicantID := uint64(ID)
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
