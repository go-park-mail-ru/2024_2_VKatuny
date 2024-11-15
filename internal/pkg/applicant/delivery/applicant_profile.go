package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio"
	"github.com/sirupsen/logrus"
)

type ApplicantHandlers struct {
	logger           *logrus.Entry
	backandURI       string
	applicantUsecase applicantUsecase.IApplicantUsecase
	portfolioUsecase portfolioUsecase.IPortfolioUsecase
	cvUsecase        cvs.ICVsUsecase
}

func NewApplicantProfileHandlers(app *internal.App) (*ApplicantHandlers, error) {

	return &ApplicantHandlers{
		logger:           logrus.NewEntry(app.Logger),
		applicantUsecase: app.Usecases.ApplicantUsecase,
		portfolioUsecase: app.Usecases.PortfolioUsecase,
		cvUsecase:        app.Usecases.CVUsecase,
	}, nil
}

func (h *ApplicantHandlers) ApplicantProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetApplicantProfileHandler(w, r)
	case http.MethodPut:
		h.UpdateApplicantProfileHandler(w, r)
	default:
		middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
			HTTPStatus: http.StatusMethodNotAllowed,
			Error:      http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

func (h *ApplicantHandlers) GetApplicantProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantProfileHandler"

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/applicant/profile/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
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

func (h *ApplicantHandlers) UpdateApplicantProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.UpdateApplicantProfileHandler"

	// url := r.URL.Path[len("/api/v1/applicant/profile/"):]
	// slugID := strings.Split(url, "/")[0]
	// ID, err := strconv.Atoi(slugID)
	// if len(url) < 1 || err != nil {
	// 	h.logger.Errorf("function %s: something bad with slug", fn)
	// 	middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
	// 		HTTPStatus: http.StatusBadRequest,
	// 		Error:      "something bad with slug",
	// 	})
	// 	return
	// }

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/applicant/profile/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
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

func (h *ApplicantHandlers) GetApplicantPortfoliosHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantPortfoliosHandler"

	// url := r.URL.Path[len("/api/v1/applicant/portfolio/"):]
	// slugID := strings.Split(url, "/")[0]
	// ID, err := strconv.Atoi(slugID)
	// if len(url) < 1 || err != nil {
	// 	h.logger.Errorf("function %s: something bad with slug", fn)
	// 	middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
	// 		HTTPStatus: http.StatusBadRequest,
	// 		Error:      "something bad with slug",
	// 	})
	// 	return
	// }
	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/applicant/portfolio/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
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

func (h *ApplicantHandlers) GetApplicantCVsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "ApplicantProfileHandlers.GetApplicantCVsHandler"

	// url := r.URL.Path[len("/api/v1/applicant/cv/"):]
	// slugID := strings.Split(url, "/")[0]
	// ID, err := strconv.Atoi(slugID)
	// if len(url) < 1 || err != nil {
	// 	h.logger.Errorf("function %s: something bad with slug", fn)
	// 	middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
	// 		HTTPStatus: http.StatusBadRequest,
	// 		Error:      "something bad with slug",
	// 	})
	// 	return
	// }
	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/applicant/cv/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
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
