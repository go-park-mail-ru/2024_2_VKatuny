// Package delivery is a handlers layer of employer
package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

// CreateEmployerHandler creates employers in db
// CreateEmployer godoc
// @Summary     Creates a new user as a employer
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Success     200      {object}       dto.JSONResponse{statusCode=200,body=dto.JSONUserBody, error=""} "OK"
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

		newUserInput := new(dto.JSONEmployerRegistrationForm)
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("error while unmarshalling employer  JSON: %s", err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}
		if err := employerUsecase.CreateEmployerInputCheck(newUserInput); err != nil {
			logger.Errorf("employer invalid fields")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		// TODO: remake
		newUserInput.Password = utils.HashPassword(newUserInput.Password)
		userID, err := repo.Create(newUserInput)
		logger.Debugf("function %s: employer created", funcName)
		if err == nil {
			// TODO: cover error
			user, _ := repo.GetByID(userID)
			middleware.UniversalMarshal(w, http.StatusOK, user)
			return
		}
		logger.Debugf("error user with this email already exists: %s", newUserInput.Email)
		middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      "user already exist",
		})
	})
}
