package delivery

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
// 	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
// 	"github.com/sirupsen/logrus"
// )

// type EmployerProfileHandlers struct {
// 	logger *logrus.Logger
// }


// func (h *EmployerProfileHandlers) GetEmployerProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	fn := "EmployerProfileHandlers.GetEmployerProfileHandler"

// 	ID, err := middleware.GetIDSlugAtEnd(w, r, "/api/v1/profile/employer")
// 	if err != nil {
// 		h.logger.Errorf("function %s: got err %s", fn, err)
// 		return
// 	}

// 	employerID := uint64(ID)
// 	// dto - JSONGetEmployerProfile
// 	employerProfile, err := h.employerUsecase.GetEmployerProfile(employerID)
// 	if err != nil {
// 		h.logger.Errorf("function %s: got err %s", fn, err)
// 		middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
// 			HTTPStatus: http.StatusInternalServerError,
// 			Error:      err.Error(),
// 		})
// 		return
// 	}

// 	h.logger.Debugf("function %s: success, got profile %v", fn, employerProfile)
// 	middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
// 		HTTPStatus: http.StatusOK,
// 		Body:       employerProfile,
// 	})
// }

// func (h *EmployerProfileHandlers) UpdateEmployerProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	fn := "EmployerProfileHandlers.UpdateEmployerProfileHandler"


