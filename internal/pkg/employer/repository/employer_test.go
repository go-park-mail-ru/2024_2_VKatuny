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
					mock.ExpectQuery(`select employer.id, first_name, last_name, city.city_name, position, company.company_name, company_description, company_website, path_to_profile_avatar, contacts, 
						email, password_hash, employer.created_at, employer.updated_at 
						from employer left join city on employer.city_id = city.id left join company on employer.company_name_id = company.id where employer.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "city.city_name",
							"position", "company.company_name", "company_description", "company_website", "path_to_profile_avatar", "contacts", "email", "password_hash", "employer.created_at", "employer.updated_at"}).
							AddRow(1, "Иван", "Иванов", "Москва", "Должность", "Компания", "Описание компании", "Сайт компании", "/src",
								"tg - ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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
					mock.ExpectQuery(`select employer.id, first_name, last_name, city.city_name, position, company.company_name, company_description, company_website, path_to_profile_avatar, contacts, 
						email, password_hash, employer.created_at, employer.updated_at 
						from employer left join city on employer.city_id = city.id left join company on employer.company_name_id = company.id where employer.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "city.city_name",
							"position", "company.company_name", "company_description", "company_website", "path_to_profile_avatar", "contacts", "email", "password_hash", "employer.created_at", "employer.updated_at"}).
							AddRow(1, nil, nil, nil, nil, nil, nil, nil, nil,
								nil, nil, nil, nil, nil))
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

			s := NewEmployerStorage(db)

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
				Email: "a.mail.ru",
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select employer.id, first_name, last_name, city.city_name, position, company.company_name, company_description, company_website, path_to_profile_avatar, contacts, 
						email, password_hash, employer.created_at, employer.updated_at 
						from employer left join city on employer.city_id = city.id left join company on employer.company_name_id = company.id where employer.email = (.+)`).
						WithArgs(
							args.Email,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "city.city_name",
							"position", "company.company_name", "company_description", "company_website", "path_to_profile_avatar", "contacts", "email", "password_hash", "employer.created_at", "employer.updated_at"}).
							AddRow(1, "Иван", "Иванов", "Москва", "Должность", "Компания", "Описание компании", "Сайт компании", "/src",
								"tg - ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewEmployerStorage(db)

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
		employer  dto.EmployerInput
		companyID uint64
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
				employer: dto.EmployerInput{
					FirstName:          "Иван",
					LastName:           "Position",
					Position:           "Должность",
					CompanyName:        "Apple",
					CompanyDescription: "Описание компании",
					CompanyWebsite:     "Вебсайт",
					Contacts:           "tg - ",
					Email:              "a@mail.ru",
					Password:           "12341234",
				},
				companyID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.employer.CompanyName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into employer (.+)`).
						WithArgs(
							args.employer.FirstName,
							args.employer.LastName,
							args.employer.Position,
							args.companyID,
							args.employer.CompanyDescription,
							args.employer.CompanyWebsite,
							args.employer.Email,
							args.employer.Password,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name",
							"position", "company_description", "company_website", "path_to_profile_avatar", "contacts", "email", "password_hash", "employer.created_at", "employer.updated_at", "employer.compressed_image"}).
							AddRow(1, "Иван", "Иванов", "Должность", "Описание компании", "Сайт компании", "/src",
								"tg - ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "compressed_image"))
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

			s := NewEmployerStorage(db)

			if _, err := s.Create(&tt.args.employer); (err != nil) != tt.wantErr {
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
		ID       uint64
		CityID   uint64
		employer dto.JSONUpdateEmployerProfile
		query1   func(mock sqlmock.Sqlmock, args args)
		query2   func(mock sqlmock.Sqlmock, args args)
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
				employer: dto.JSONUpdateEmployerProfile{
					FirstName: "Иван",
					LastName:  "Иванов",
					City:      "Москва",
					Contacts:  "tg - ",
				},
				CityID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.employer.City,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update employer (.+)`).
						WithArgs(
							args.employer.FirstName,
							args.employer.LastName,
							args.CityID,
							args.employer.Contacts,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name",
							"position", "company_description", "company_website", "path_to_profile_avatar", "contacts", "email", "password_hash", "employer.created_at", "employer.updated_at", "employer.compressed_image"}).
							AddRow(1, "Иван", "Иванов", "Должность", "Описание компании", "Сайт компании", "/src",
								"tg - ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "compressed_image"))
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

			s := NewEmployerStorage(db)

			if _, err := s.Update(tt.args.ID, &tt.args.employer); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
