package notificationmicroservice

import (
	"context"

	notifications "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"
	notificationsinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications"
	"github.com/sirupsen/logrus"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
)

type NotificationsManager struct {
	notifications.UnsafeNotificationsServiceServer
	notificationsUsecase notificationsinterfaces.INotificationsUsecase
	logger               *logrus.Entry
}

func NewNotificationsManager(notificationsUsecase notificationsinterfaces.INotificationsUsecase, logger *logrus.Logger) *NotificationsManager {
	return &NotificationsManager{
		notificationsUsecase: notificationsUsecase,
		logger:              logrus.NewEntry(logger),
	}
}

// func (nd *NotificationsManager) GetAlEmployerNotifications(ctx context.Context, in *notifications.GetAlEmployerNotificationsInput) (*notifications.GetAlEmployerNotificationsOutput, error) {
// 	funcName := "NotificationsDelivery.GetAlEmployerNotifications"
// 	nd.logger.Debugf("%s: got request: %s", funcName, in)
// 	if in == nil {
// 		return nil, notificationsinterfaces.ErrNothingInInputData
// 	}
// 	notificationsList, err := nd.notificationsUsecase.GetAlEmployerNotifications(in.EmployerID)
// 	out := &notifications.GetAlEmployerNotificationsOutput{}
// 	for _, oneNotification := range notificationsList {
// 		out.Notifications = append(out.Notifications, &notifications.Notification{
// 			ID:               oneNotification.ID,
// 			NotificationText: oneNotification.NotificationText,
// 			Read:             oneNotification.IsRead,
// 			CreatedAt:        oneNotification.CreatedAt,
// 		})
// 	}
// 	return out, err
// }

// func (nd *NotificationsManager) MakeEmployerNotificationRead(ctx context.Context, in *notifications.MakeEmployerNotificationReadInput) (*notifications.Nothing, error) {
// 	funcName := "NotificationsDelivery.MakeEmployerNotificationRead"
// 	nd.logger.Debugf("%s: got request: %s", funcName, in)
// 	if in == nil {
// 		return &notifications.Nothing{}, notificationsinterfaces.ErrNothingInInputData
// 	}
// 	err := nd.notificationsUsecase.MakeEmployerNotificationRead(in.NotificationID)
// 	return &notifications.Nothing{}, err
// }

func (nd *NotificationsManager) CreateEmployerNotification(ctx context.Context, in *notifications.CreateEmployerNotificationInput) (*notifications.Nothing, error) {
	funcName := "NotificationsDelivery.CreateEmployerNotification"
	nd.logger.Debugf("%s: got request: %s", funcName, in)
	if in == nil {
		return &notifications.Nothing{}, dto.ErrNothingInInputData
	}
	err := nd.notificationsUsecase.CreateEmployerNotification(in.ApplicantID, in.EmployerID, in.VacancyID, in.ApplicantInfo, in.VacancyInfo)
	return &notifications.Nothing{}, err
}
