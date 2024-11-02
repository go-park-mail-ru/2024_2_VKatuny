package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker"
	workerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker/usecase"
	"github.com/sirupsen/logrus"
)

// CreateWorkerHandler creates worker in db
// CreateWorker godoc
// @Summary     Creates a new user as a worker
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} inmemorydb.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/worker/ [post]
func CreateWorkerHandler(repo worker.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "CreateWorkerHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
			return
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(models.Worker)
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("function %s: got err %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}

		if err := workerUsecase.CreateWorkerInputCheck(newUserInput.Name, newUserInput.LastName, newUserInput.Email, newUserInput.Password); err != nil {
			logger.Errorf("function %s: %s", funcName, err.Error())
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "user's fields aren't valid",
			})
			return
		}

		logger.Debugf("function %s: adding applicant to db %v", funcName, newUserInput)
		userID, err := repo.Add(newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Body: dto.JSONUserBody{
					UserType: "applicant",
					ID:       userID,
				}, // check in postman
			})
		} else {
			// is there actually should be HTTP 400?
			logger.Errorf("function %s: got err while adding applicant to db %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't add applicant to db",
			})
		}

	})
}
