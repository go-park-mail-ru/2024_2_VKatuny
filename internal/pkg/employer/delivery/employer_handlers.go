// Package delivery is a handlers layer of employer
package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	sessionRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
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
func CreateEmployerHandler(repo repository.EmployerRepository, repoEmployerSession sessionRepo.SessionRepository, backendAddress string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "CreateEmployerHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
			return
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(dto.EmployerInput)
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
		logger.Debugf("function %s: employer input check passed", funcName)

		user, sessionID, err := employerUsecase.CreateEmployer(repo, repoEmployerSession, newUserInput)
		if err != nil {
			logger.Errorf("employer invalid fields")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      err.Error(),
			})
			return
		}
		logger.Debug("Cookie send")
		cookie := utils.MakeAuthCookie(sessionID, backendAddress)
		http.SetCookie(w, cookie)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Body:       user,
			})
		} else {
			// is there actually should be HTTP 400?
			logger.Errorf("function %s: got err while adding applicant to db %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
	})
}
