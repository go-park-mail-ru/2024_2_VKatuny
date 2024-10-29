package workerDelivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker"
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
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(models.Worker)
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("function %s: got err %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}
		if len(newUserInput.Name) < 3 || len(newUserInput.LastName) < 3 ||
			strings.Index(newUserInput.Email, "@") < 0 || len(newUserInput.Password) < 4 {
			logger.Errorf("function %s: User's fields aren't valid %+v", funcName, newUserInput)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "user's fields aren't valid",
			})
			return
		}
		logger.Debugf("function %s: adding applicant to db %v", funcName, newUserInput)
		userId, err := repo.Add(newUserInput)
		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, dto.JsonResponse{
				HttpStatus: http.StatusOK,
				Body: dto.JsonUserBody{
					UserType: "applicant",
					ID:       userId,
				}, // check in postman
			})
		} else {
			// is there actually should be HTTP 400?
			logger.Errorf("function %s: got err while adding applicant to db %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JsonResponse{
				HttpStatus: http.StatusBadRequest,
				Error:      "can't add applicant to db",
			})
		}

	})
}
