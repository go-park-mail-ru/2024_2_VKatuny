package dto

import "fmt"

type EmployerNotification struct {
	ID               uint64
	NotificationText string
	ApplicantID      uint64
	EmployerID       uint64
	VacancyID        uint64
	IsRead           bool
	CreatedAt        string
}

// JSONResponse is a standard form of response from backend to frontend
type JSONResponse struct {
	HTTPStatus int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Error      string      `json:"error"`
}

var (
	ErrNothingInInputData = fmt.Errorf("Nothing in input data")
	MsgInvalidJSON        = "invalid json"
)
