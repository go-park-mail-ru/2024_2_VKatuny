package survey

import "fmt"

var (
	ErrInvalidJSON = fmt.Errorf("invalid JSON")
)

type JSONSurveyStatistics struct{}

type JSONSurveyForm struct{}

type ISurveyUsecase interface {
	GetStatistics() (*JSONSurveyStatistics, error)
	AddAnswer(*JSONSurveyForm, string) error
	GetForm(surveyType string) (*JSONSurveyForm, error)
}

