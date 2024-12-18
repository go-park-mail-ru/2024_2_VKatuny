package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
)

func TestGetAlEmployerNotifications(t *testing.T) {
	t.Parallel()
	type args struct {
		employerID    uint64
		query func(mock sqlmock.Sqlmock, args args)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "TestOk",
			args: args{
				employerID: 1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select id, notification_text, applicant_id, employer_id, vacancy_id, is_read, created_at
								from employer_notification where employer_notification.employer_id = (.+)`).
						WithArgs(
							args.employerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "notification_text", "applicant_id", "employer_id","vacancy_id", "is_read", "created_at"}).
							AddRow(1, "Notification text", 1, 1, 1, true, "2024-11-09 04:17:52.598 +0300"))
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.args.query(mock, tt.args)

			logger := logger.NewLogrusLogger()
			s := NewNotificationsRepository(logger, db)

			if _, err := s.GetAlEmployerNotifications(tt.args.employerID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v", err != nil, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCreateEmployerNotification(t *testing.T) {
	t.Parallel()
	type args struct {
		applicantID      uint64
		employerID       uint64
		vacancyID        uint64
		NotificationText string
		query1           func(mock sqlmock.Sqlmock, args args)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "TestOk",
			args: args{
				applicantID:      1,
				employerID:       1,
				vacancyID:        1,
				NotificationText: "text",
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectExec(`insert into employer_notification (.+)`).
						WithArgs(
							args.applicantID,
							args.employerID,
							args.vacancyID,
							args.NotificationText,
						).WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.args.query1(mock, tt.args)
			logger := logger.NewLogrusLogger()
			s := NewNotificationsRepository(logger, db)

			if err := s.CreateEmployerNotification(tt.args.applicantID, tt.args.employerID, tt.args.vacancyID, tt.args.NotificationText); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestMakeEmployerNotificationRead(t *testing.T) {
	t.Parallel()
	type args struct {
		notificationID    uint64
		query1    func(mock sqlmock.Sqlmock, args args)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "TestOk",
			args: args{
				notificationID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update employer_notification (.+)`).
						WithArgs(
							args.notificationID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.args.query1(mock, tt.args)

			logger := logger.NewLogrusLogger()
			s := NewNotificationsRepository(logger, db)

			if err := s.MakeEmployerNotificationRead(tt.args.notificationID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
