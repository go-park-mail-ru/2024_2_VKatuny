package dto

import "fmt"

type EmployerNotification struct {
	ID               uint64 `json:"id"`
	NotificationText string `json:"notificationText"`
	ApplicantID      uint64 `json:"applicantId"`
	EmployerID       uint64 `json:"employerId"`
	VacancyID        uint64 `json:"vacancyId"`
	IsRead           bool   `json:"isRead"`
	CreatedAt        string `json:"createdAt"`
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
