package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
)

// func TestPostgresGetByID(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ID    uint64
// 		query func(mock sqlmock.Sqlmock, args args)
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 		err     error
// 	}{
// 		{
// 			name: "TestOk",
// 			args: args{
// 				ID: 1,
// 				query: func(mock sqlmock.Sqlmock, args args) {
// 					mock.ExpectQuery(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts,
// 						education, email, password_hash, applicant.created_at, applicant.updated_at, applicant.compressed_image
// 						from applicant left join city on applicant.city_id = city.id where applicant.id = (.+)`).
// 						WithArgs(
// 							args.ID,
// 						).
// 						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
// 							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at", "compressed_image"}).
// 							AddRow(1, "Иван", "Иванов", "Москва", "12-12-2001", "/src",
// 								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "image"))
// 				},
// 			},
// 			wantErr: false,
// 			err:     nil,
// 		},
// 		{
// 			name: "TestFailZeroID",
// 			args: args{
// 				ID: 0,
// 				query: func(mock sqlmock.Sqlmock, args args) {
// 					mock.ExpectQuery(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts,
// 						education, email, password_hash, applicant.created_at, applicant.updated_at, applicant.compressed_image
// 						from applicant left join city on applicant.city_id = city.id where applicant.id = (.+)`).
// 						WithArgs(
// 							args.ID,
// 						).
// 						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
// 							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at", "compressed_image"}).
// 							AddRow(1, nil, nil, nil, nil, nil,
// 								nil, nil, nil, nil, nil, nil, nil))
// 				},
// 			},
// 			wantErr: true,
// 			err:     nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			db, mock, err := sqlmock.New()
// 			if err != nil {
// 				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 			}
// 			defer db.Close()

// 			tt.args.query(mock, tt.args)

// 			s := NewNotificationsRepository(db)

// 			if _, err := s.GetAlEmployerNotifications(tt.args.ID); (err != nil) != tt.wantErr {
// 				t.Errorf("Postgres error = %v, wantErr %v, err!!! %s", err != nil, tt.wantErr, err)
// 			}

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("there were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

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

// func TestMakeEmployerNotificationRead(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ID        uint64
// 		CityID    uint64
// 		applicant dto.JSONUpdateApplicantProfile
// 		query1    func(mock sqlmock.Sqlmock, args args)
// 		query2    func(mock sqlmock.Sqlmock, args args)
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 		err     error
// 	}{
// 		{
// 			name: "TestOk",
// 			args: args{
// 				ID: 1,
// 				applicant: dto.JSONUpdateApplicantProfile{
// 					FirstName: "Иван",
// 					LastName:  "Иванов",
// 					City:      "Москва",
// 					BirthDate: "12-12-2001",
// 					Contacts:  "tg - ",
// 					Education: "Высшее",
// 				},
// 				CityID: 1,
// 				query1: func(mock sqlmock.Sqlmock, args args) {
// 					mock.ExpectQuery(`select  (.+)`).
// 						WithArgs(
// 							args.applicant.City,
// 						).
// 						WillReturnRows(sqlmock.NewRows([]string{"id"}).
// 							AddRow(1))
// 				},
// 				query2: func(mock sqlmock.Sqlmock, args args) {
// 					mock.ExpectQuery(`update applicant (.+)`).
// 						WithArgs(
// 							args.applicant.FirstName,
// 							args.applicant.LastName,
// 							args.CityID,
// 							args.applicant.BirthDate,
// 							args.applicant.Contacts,
// 							args.applicant.Education,
// 							args.ID,
// 						).
// 						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name",
// 							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at", "compressed_image"}).
// 							AddRow(1, "Иван", "Иванов", "12-12-2001", "/src",
// 								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "image"))
// 				},
// 			},
// 			wantErr: false,
// 			err:     nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			db, mock, err := sqlmock.New()
// 			if err != nil {
// 				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 			}
// 			defer db.Close()

// 			tt.args.query1(mock, tt.args)
// 			tt.args.query2(mock, tt.args)

// 			logger := logger.NewLogrusLogger()
// 			s := NewNotificationsRepository(logger, db)

// 			if err := s.MakeEmployerNotificationRead(tt.args.ID); (err != nil) != tt.wantErr {
// 				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
// 			}

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("there were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }
