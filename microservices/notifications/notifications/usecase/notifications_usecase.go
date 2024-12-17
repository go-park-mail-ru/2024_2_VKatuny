package usecase

import (
	notificationsinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
	"github.com/sirupsen/logrus"
)

type NotificationsUsecase struct {
	notificationsRepo notificationsinterfaces.INotificationsRepository
	logger            *logrus.Entry
}

func NewNotificationsUsecase(notificationsRepo notificationsinterfaces.INotificationsRepository, logger *logrus.Logger) *NotificationsUsecase {
	return &NotificationsUsecase{
		notificationsRepo: notificationsRepo,
		logger:            &logrus.Entry{Logger: logger},
	}
}

func (cu *NotificationsUsecase) GetAlEmployerNotifications(employerID uint64) ([]*dto.EmployerNotification, error) {
	funcName := "NotificationsUsecase.GetAlEmployerNotifications"
	cu.logger.Debugf("%s: got request: %d", funcName, employerID)
	notificationsList, err := cu.notificationsRepo.GetAlEmployerNotifications(employerID)
	return notificationsList, err
}

func (cu *NotificationsUsecase) MakeEmployerNotificationRead(notificationID uint64) error {
	funcName := "NotificationsUsecase.MakeEmployerNotificationRead"
	cu.logger.Debugf("%s: got request: %d", funcName, notificationID)
	err := cu.notificationsRepo.MakeEmployerNotificationRead(notificationID)
	return err
}

func (cu *NotificationsUsecase) CreateEmployerNotification(applicantID, employerID, vacancyID uint64, applicantInfo, vacancyInfo string) error {
	funcName := "NotificationsUsecase.CreateEmployerNotification"
	cu.logger.Debugf("%s: got request: %d %d %d %s %s", funcName, applicantID, employerID, vacancyID, applicantInfo, vacancyInfo)
	notificationText := `На вашу вакансию "`+vacancyInfo+`" откликнулся ` + applicantInfo
	err := cu.notificationsRepo.CreateEmployerNotification(applicantID, employerID, vacancyID, notificationText)
	return err
}
