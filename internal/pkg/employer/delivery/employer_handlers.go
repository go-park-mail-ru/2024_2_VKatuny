package employerDelivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer"
	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/sirupsen/logrus"
)

// CreateEmployerHandler creates employers in db
// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string    true         "User's email"
// @Param       password body string    true         "User's password"
// @Success     200      {object}       inmemorydb.UserInput
// @Failure     400      {object}       nil
// @Router      /registration/employer/ [post]
func CreateEmployerHandler(repo employer.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "CreateEmployerHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
			return
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(models.Employer)
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("error while unmarshalling employer  JSON: %s", err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}
		if err := employerUsecase.CreateEmployerInputCheck(newUserInput.Name, newUserInput.LastName, newUserInput.Position, newUserInput.CompanyName, newUserInput.Email, newUserInput.Password); err != nil {
			logger.Errorf("employer invalid fields")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		user, err := repo.Create(newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, user)
			return
		}
		logger.Debugf("error user with this email already exists: %s", newUserInput.Email)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      "user already exist",
		})
	})
}
