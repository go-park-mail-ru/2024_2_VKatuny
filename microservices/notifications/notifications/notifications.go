package notificationsmicroserviceinterface

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
)

// Interface for Notifications.
type INotificationsRepository interface {
	GetAlEmployerNotifications(employerID uint64) ([]*dto.EmployerNotification, error)
	MakeEmployerNotificationRead(notificationID uint64) error
	CreateEmployerNotification(applicantID uint64, employerID uint64, vacancyID uint64, NotificationText string) error
}

type INotificationsUsecase interface {
	GetAlEmployerNotifications(employerID uint64) ([]*dto.EmployerNotification, error)
	MakeEmployerNotificationRead(notificationID uint64) error
	CreateEmployerNotification(applicantID, employerID, vacancyID uint64, applicantInfo, vacancyInfo string) error
}
