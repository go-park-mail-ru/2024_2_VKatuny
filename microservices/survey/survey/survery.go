package survey

import "fmt"

// Interface for Compress.
type ISurveryRepository interface {
	GetStatistic() ([]*Statistics, error)
	GetQuestionByType() ([]*Question, error)
	CreateAnswerAuthorised(QuestionAnswer *QuestionAnswer) error
}

type ISurveyUsecase interface {
	GetStatistics() (*JSONSurveyStatistics, error)
	AddAnswer(*JSONSurveyForm, string) error
	GetForm(surveyType string) (*JSONSurveyForm, error)
}

var (
	ErrInvalidJSON = fmt.Errorf("invalid JSON")
	ErrUnableToGetStatistics = fmt.Errorf()
)

type JSONSurveyStatistics struct{
	ValAVG       int32 `json:"avgRating"`
	QuestionText int32 `json:"questionText"`
	QuestionID   int32 `json:"questionID"`
}

type JSONSurveyForm struct{}

type QuestionAnswer struct {
	QuestionID int32  `json:"questionID"`
	Token      string `json:"token"`
	Value      int32  `json:"value"`
}

type Question struct {
	ID           int32  `json:"ID"`
	QuestionText string `json:"question_text"`
	TypeText     string `json:"typeText"`
	Position     int32  `json:"position"`
}

type Statistics struct {
	ValAVG       int32 `json:"avgRating"`
	QuestionText int32 `json:"questionText"`
	QuestionID   int32 `json:"questionID"`
}

type JSONResponse struct {
	HTTPStatus int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Error      string      `json:"error"`
}
