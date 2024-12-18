package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
	"github.com/sirupsen/logrus"
)

type NotificationsRepository struct {
	logger *logrus.Entry
	db     *sql.DB
}

func NewNotificationsRepository(logger *logrus.Logger, db *sql.DB) *NotificationsRepository {
	return &NotificationsRepository{
		logger: &logrus.Entry{Logger: logger},
		db:     db,
	}
}

func (nr *NotificationsRepository) GetAlEmployerNotifications(employerID uint64) ([]*dto.EmployerNotification, error) {
	funcName := "NotificationsRepository.GetAlEmployerNotifications"
	nr.logger.Debugf("%s: got request: %d", funcName, employerID)
	NotificationsList := make([]*dto.EmployerNotification, 0)

	rows, err := nr.db.Query(`select id, notification_text, applicant_id, employer_id, vacancy_id, is_read, created_at
	from employer_notification where employer_notification.employer_id = $1`, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneNotification dto.EmployerNotification
		if err := rows.Scan(
			&oneNotification.ID,
			&oneNotification.NotificationText,
			&oneNotification.ApplicantID,
			&oneNotification.EmployerID,
			&oneNotification.VacancyID,
			&oneNotification.IsRead,
			&oneNotification.CreatedAt,
		); err != nil {
			return nil, err
		}
		NotificationsList = append(NotificationsList, &oneNotification)
		nr.logger.Debugf("notification: %+v", oneNotification)
	}
	return NotificationsList, nil
}

func (nr *NotificationsRepository) MakeEmployerNotificationRead(notificationID uint64) error {
	funcName := "NotificationsRepository.MakeEmployerNotificationRead"
	nr.logger.Debugf("%s: got request: %d", funcName, notificationID)
	row := nr.db.QueryRow(`update employer_notification set is_read = true where id = $1 returning id`, notificationID)
	var id uint64
	err := row.Scan(
		&id)
	return err
}

func (nr *NotificationsRepository) CreateEmployerNotification(applicantID uint64, employerID uint64, vacancyID uint64, NotificationText string) error {
	funcName := "NotificationsRepository.CreateEmployerNotification"
	nr.logger.Debugf("%s: got request: %d %d %d %s", funcName, applicantID, employerID, vacancyID, NotificationText)
	_, err := nr.db.Exec(`insert into employer_notification (applicant_id, employer_id, vacancy_id, notification_text) VALUES ($1, $2, $3, $4) returning id`, applicantID, employerID, vacancyID, NotificationText)
	return err
}
