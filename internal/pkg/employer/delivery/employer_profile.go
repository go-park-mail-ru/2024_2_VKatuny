package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	vacanciesUsecase"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"
	"github.com/sirupsen/logrus"
)

type EmployerProfileHandlers struct {
	logger           *logrus.Logger
	employerUsecase  *employerUsecase.EmployerUsecase
	vacanciesUsecase *vacanciesUsecase.VacanciesUsecase
}

func NewEmployerProfileHandlers(logger *logrus.Logger, usecases *internal.Usecases) (*EmployerProfileHandlers, error) {
	employerUsecase, ok1 := usecases.EmployerUsecase.(*employerUsecase.EmployerUsecase)
	vacanciesUsecase, ok2 := usecases.VacanciesUsecase.(*vacanciesUsecase.VacanciesUsecase)
	if !(ok1 && ok2) {
		return nil, commonerrors.ErrUnableToCast
	}
	return &EmployerProfileHandlers{
		logger:           logger,
		employerUsecase:  employerUsecase,
		vacanciesUsecase: vacanciesUsecase,
	}, nil
}

func (h *EmployerProfileHandlers) EmployerProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/employer/profile/" {
		switch r.Method {
		case http.MethodGet:
			h.GetEmployerProfileHandler(w, r)
		case http.MethodPut:
			h.UpdateEmployerProfileHandler(w, r)
		default:
			middleware.UniversalMarshal(w, http.StatusMethodNotAllowed, dto.JSONResponse{
				HTTPStatus: http.StatusMethodNotAllowed,
				Error:      http.StatusText(http.StatusMethodNotAllowed),
			})
		}
	}
}

func (h *EmployerProfileHandlers) GetEmployerProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerProfileHandler"

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/employer/profile/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}

	employerID := uint64(ID)
	// dto - JSONGetEmployerProfile
	employerProfile, err := h.employerUsecase.GetEmployerProfile(employerID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got profile %v", fn, employerProfile)
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       employerProfile,
	})
}

func (h *EmployerProfileHandlers) UpdateEmployerProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.UpdateEmployerProfileHandler"
	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/employer/profile/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}

	employerID := uint64(ID)

	decoder := json.NewDecoder(r.Body)
	newProfileData := new(dto.JSONUpdateEmployerProfile)
	err = decoder.Decode(newProfileData)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "unable to unmarshal JSON",
		})
		return
	}

	h.logger.Debugf("function %s: new profile data JSON parsed: %v", fn, newProfileData)

	err = h.employerUsecase.UpdateEmployerProfile(employerID, newProfileData)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("function %s: success", fn)
}

func (h *EmployerProfileHandlers) GetEmployerVacanciesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerVacanciesHandler"

	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/employer/vacancies/")
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		return
	}

	employerID := uint64(ID)
	vacancies, err := h.vacanciesUsecase.GetVacanciesByEmployerID(employerID)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	h.logger.Debugf("function %s: success, got vacancies: %d", fn, len(vacancies))
	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       vacancies,
	})
}
