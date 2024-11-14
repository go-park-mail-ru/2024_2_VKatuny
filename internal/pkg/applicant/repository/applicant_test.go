package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

func TestPostgresGetByID(t *testing.T) {
	t.Parallel()
	type args struct {
		ID    uint64
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
				ID: 1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts, 
						education, email, password_hash, applicant.created_at, applicant.updated_at 
						from applicant left join city on applicant.city_id = city.id where applicant.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at"}).
							AddRow(1, "Иван", "Иванов", "Москва", "12-12-2001", "/src",
								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestFailZeroID",
			args: args{
				ID: 0,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts, 
						education, email, password_hash, applicant.created_at, applicant.updated_at 
						from applicant left join city on applicant.city_id = city.id where applicant.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at"}).
							AddRow(1, nil, nil, nil, nil, nil,
								nil, nil, nil, nil, nil, nil))
				},
			},
			wantErr: true,
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

			s := NewApplicantStorage(db)

			if _, err := s.GetByID(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!! %s", err != nil, tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresGetByEmail(t *testing.T) {
	t.Parallel()
	type args struct {
		Email string
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
				Email: "a@mail.ru",
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select applicant.id, first_name, last_name, city.city_name, birth_date, path_to_profile_avatar, contacts, 
						education, email, password_hash, applicant.created_at, applicant.updated_at 
						from applicant left join city on applicant.city_id = city.id where applicant.email=(.+)`).
						WithArgs(
							args.Email,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at"}).
							AddRow(1, "Иван", "Иванов", "Москва", "12-12-2001", "/src",
								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewApplicantStorage(db)

			if _, err := s.GetByEmail(tt.args.Email); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!! %s", err != nil, tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresCreate(t *testing.T) {
	t.Parallel()
	type args struct {
		applicant dto.ApplicantInput
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
				applicant: dto.ApplicantInput{
					FirstName: "Иван",
					LastName:  "Position",
					BirthDate: "12-12-2001",
					Education: "Высшее",
					Email:     "a@mail.ru",
					Password:  "12341234",
				},
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into applicant (.+)`).
						WithArgs(
							args.applicant.FirstName,
							args.applicant.LastName,
							args.applicant.BirthDate,
							args.applicant.Education,
							args.applicant.Email,
							args.applicant.Password,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at"}).
							AddRow(1, "Иван", "Иванов", "12-12-2001", "/src",
								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewApplicantStorage(db)

			if _, err := s.Create(&tt.args.applicant); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		ID        uint64
		CityID    uint64
		applicant dto.JSONUpdateApplicantProfile
		query1    func(mock sqlmock.Sqlmock, args args)
		query2    func(mock sqlmock.Sqlmock, args args)
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
				ID: 1,
				applicant: dto.JSONUpdateApplicantProfile{
					FirstName: "Иван",
					LastName:  "Иванов",
					City:      "Москва",
					BirthDate: "12-12-2001",
					Contacts:  "tg - ",
					Education: "Высшее",
				},
				CityID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.applicant.City,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update applicant (.+)`).
						WithArgs(
							args.applicant.FirstName,
							args.applicant.LastName,
							args.CityID,
							args.applicant.BirthDate,
							args.applicant.Contacts,
							args.applicant.Education,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at"}).
							AddRow(1, "Иван", "Иванов", "12-12-2001", "/src",
								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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
			tt.args.query2(mock, tt.args)

			s := NewApplicantStorage(db)

			if _, err := s.Update(tt.args.ID, &tt.args.applicant); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
