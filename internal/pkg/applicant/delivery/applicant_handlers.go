package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	sessionRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

// CreateWorkerHandler creates applicant in db
// CreateWorker godoc
// @Summary     Creates a new user as a applicant
// @Description -
// @Tags        Registration
// @Accept      json
// @Produce     json
// @Param       email    body string true "User's email"
// @Param       password body string true "User's password"
// @Success     200 {object} inmemorydb.UserInput
// @Failure     http.StatusBadRequest {object} nil
// @Router      /registration/applicant/ [post]
func CreateApplicantHandler(repo repository.ApplicantRepository, repoApplicantSession sessionRepo.SessionRepository, backendAddress string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		funcName := "CreateApplicantHandler"
		logger, ok := r.Context().Value(dto.LoggerContextKey).(*logrus.Logger)
		if !ok {
			fmt.Printf("function %s: can't get logger from context\n", funcName)
			return
		}

		decoder := json.NewDecoder(r.Body)

		newUserInput := new(dto.ApplicantInput)
		err := decoder.Decode(newUserInput)
		if err != nil {
			logger.Errorf("function %s: got err %s", funcName, err)
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "can't unmarshal JSON",
			})
			return
		}

		if err := applicantUsecase.CreateApplicantInputCheck(newUserInput.FirstName, newUserInput.LastName, newUserInput.Email, newUserInput.Password); err != nil {
			logger.Errorf("function %s: %s", funcName, err.Error())
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusBadRequest,
				Error:      "user's fields aren't valid",
			})
			return
		}

		logger.Debugf("function %s: adding applicant to db %v", funcName, newUserInput)

		user, sessionID, err := applicantUsecase.CreateApplicant(repo, repoApplicantSession, newUserInput)
		if err != nil {
			logger.Errorf("applicant invalid fields")
			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
				HTTPStatus: http.StatusInternalServerError,
				Error:      "db err",
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
