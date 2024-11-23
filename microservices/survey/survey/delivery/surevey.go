package delivery

import (
	"encoding/json"

	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/survey/survey"

	"github.com/sirupsen/logrus"
)

type SurveyHandlers struct {
	logger        *logrus.Entry
	surveyUsecase survey.ISurveyUsecase
	getParams     map[string]string
}

func NewSurveyHandlers(logger *logrus.Logger) *SurveyHandlers {
	return &SurveyHandlers{
		logger: logrus.NewEntry(logger),
		getParams: map[string]string{
			"csat": "------", // TODO: set get parameters
		},
	}
}

func (h *SurveyHandlers) GetStatistics(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SurveyHandlers.GetStatistics"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	h.logger.Debugf("%s: entering", fn)

	// TODO: implement auth check

	h.logger.Debugf("%s: authorized", fn)

	statistics, err := h.surveyUsecase.GetStatistics()
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, survey.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	middleware.UniversalMarshal(w, http.StatusOK, survey.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       statistics,
	})
}

func (h *SurveyHandlers) AddSurveyAnswer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SurveyHandlers.AddSurvey"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	h.logger.Debugf("%s: entering", fn)

	surveyForm := new(survey.JSONSurveyForm)
	err := json.NewDecoder(r.Body).Decode(surveyForm)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusBadRequest, survey.JSONResponse{
			HTTPStatus: http.StatusBadRequest,
			Error:      survey.ErrInvalidJSON.Error(),
		})
		return
	}
	h.logger.Debugf("%s: survey form parsed: %v", fn, surveyForm)

	// TODO: get session
	var SessionID string

	err = h.surveyUsecase.AddAnswer(surveyForm, SessionID)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, survey.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: answer successfully added", fn)

	middleware.UniversalMarshal(w, http.StatusOK, survey.JSONResponse{
		HTTPStatus: http.StatusOK,
	})
}

func (h *SurveyHandlers) GetSurveyForm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fn := "SurveyHandlers.GetSurveyForm"
	h.logger = utils.SetRequestIDInLoggerFromRequest(r, h.logger)

	h.logger.Debugf("%s: entering", fn)

	queryParams := r.URL.Query()
	h.logger.Debugf("%s: query params: %v", fn, queryParams)

	formType := queryParams.Get(h.getParams["csat"])
	form, err := h.surveyUsecase.GetForm(formType)
	if err != nil {
		h.logger.Errorf("%s: got err %s", fn, err)
		middleware.UniversalMarshal(w, http.StatusInternalServerError, survey.JSONResponse{
			HTTPStatus: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	h.logger.Debugf("%s: got form: %v", fn, form)

	middleware.UniversalMarshal(w, http.StatusOK, survey.JSONResponse{
		HTTPStatus: http.StatusOK,
		Body:       form,
	})
}
