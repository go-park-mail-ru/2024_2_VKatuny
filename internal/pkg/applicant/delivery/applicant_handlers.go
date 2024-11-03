package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
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
func CreateApplicantHandler(repo repository.ApplicantRepository) http.Handler {
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

		user, err := applicantUsecase.CreateApplicant(repo, newUserInput)

		if err == nil {
			middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
				HTTPStatus: http.StatusOK,
				Body: dto.ApplicantOutput{
					ID:                  user.ID,
					FirstName:           user.FirstName,
					LastName:            user.LastName,
					CityName:            user.CityName,
					BirthDate:           user.BirthDate,
					PathToProfileAvatar: user.LastName,
					Constants:           user.Contacts,
					Education:           user.Education,
					Email:               user.Email,
					CreatedAt:           user.CreatedAt,
					UpdatedAt:           user.UpdatedAt,
				},
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
