package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/survey/survey"
)

type SurveyUsecase struct {
	logger *logrus.NewEntry
	surveyRepository survey.ISurveryRepository
	// TODO: add new fields
}

func NewSurveyUsecase(logger *logrus.Logger) *SurveyUsecase {
	return &SurveyUsecase{
		logger: logrus.NewEntry(logger)
		surveyRepository: 
		// TODO: add new fields
	}
}

func (u *ServeyUsecase) GetStatistics() (*JSONSurveyStatistics, error) {
	fn := "SurveyUsecase.GetStatistics"

	h.logger.Debugf("%s: entering", fn)

	statistics, err := u.surveyRepository.GetStatistics()
	if err != nil {
		u.logger.Errorf("%s: got %s", err)
		return nil, ErrUnableToGetStatistics
	}
	u.logger.Debugf("%s: statistics got without errors, len = %d", len(statistics))

	
}