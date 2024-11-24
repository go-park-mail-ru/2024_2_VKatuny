package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
)

type EmployerHandlers struct {
	logger             *logrus.Entry
	backendAddress     string
	employerUsecase    employer.IEmployerUsecase
	vacanciesUsecase   vacancies.IVacanciesUsecase
	sessionUsecase     session.ISessionUsecase
	fileLoadingUsecase fileloading.IFileLoadingUsecase
}

func NewEmployerHandlers(app *internal.App) *EmployerHandlers {
	return &EmployerHandlers{
		logger:             logrus.NewEntry(app.Logger),
		backendAddress:     app.BackendAddress,
		employerUsecase:    app.Usecases.EmployerUsecase,
		vacanciesUsecase:   app.Usecases.VacanciesUsecase,
		sessionUsecase:     app.Usecases.SessionUsecase,
		fileLoadingUsecase: app.Usecases.FileLoadingUsecase,
	}
}

func (h *EmployerHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerProfileHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)

	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}
	// dto - JSONGetEmployerProfile
	employerProfile, err := h.employerUsecase.GetEmployerProfile(r.Context(), employerID)
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

func (h *EmployerHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.UpdateEmployerProfileHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)
	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

	newProfileData := &dto.JSONUpdateEmployerProfile{}
	newProfileData.FirstName = r.FormValue("firstName")
	newProfileData.LastName = r.FormValue("lastName")
	newProfileData.City = r.FormValue("city")
	newProfileData.Contacts = r.FormValue("contacts")
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
		newProfileData.Avatar = fileAddress
	}

	h.logger.Debugf("function %s: new profile data MultiPart parsed: %v", fn, newProfileData)

	err = h.employerUsecase.UpdateEmployerProfile(r.Context(), employerID, newProfileData)
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

func (h *EmployerHandlers) GetEmployerVacancies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "EmployerProfileHandlers.GetEmployerVacanciesHandler"
	h.logger = utils.SetLoggerRequestID(r.Context(), h.logger)

	vars := mux.Vars(r)
	slug := vars["id"]
	employerID, err := strconv.ParseUint(slug, 10, 64)
	if err != nil {
		h.logger.Errorf("function %s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      commonerrors.ErrFrontUnableToCastSlug.Error(),
		})
		return
	}

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
