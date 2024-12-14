package dto

type EmployerNotification struct {
	ID               uint64
	NotificationText string
	ApplicantID      uint64
	EmployerID       uint64
	VacancyID        uint64
	IsRead             bool
	CreatedAt        string
}
