package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

func TestPostgresGetVacanciesByEmployerID(t *testing.T) {
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
					mock.ExpectQuery(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
							company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
							left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
							left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
							left join position_category on vacancy.position_category_id = position_category.id
							where vacancy.employer_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "city_name", "position", "vacancy_description",
							"salary", "employer_id", "work_type_name", "path_to_company_avatar", "created_at", "updated_at", "company_name", "category_name", "compressed_image"}).
							AddRow(1, "Moscow", "Скульптор", "Требуется скульптор без опыта работы", 90000, 1,
								"Полная занятость", "", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "Мэрия Москвы", "Творец", "test"))
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
					mock.ExpectQuery(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
							company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
							left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
							left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
							left join position_category on vacancy.position_category_id = position_category.id
							where vacancy.employer_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "city_name", "position", "vacancy_description",
							"salary", "employer_id", "work_type_name", "path_to_company_avatar", "created_at", "updated_at", "company_name", "category_name", "compressed_image"}).
							AddRow(1, nil, nil, nil, nil, nil,
								nil, nil, nil, nil, nil, nil, nil))
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

			s := NewVacanciesStorage(db)

			if _, err := s.GetVacanciesByEmployerID(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v", err != nil, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

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
					mock.ExpectQuery(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
							company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
							left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
							left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
							left join position_category on vacancy.position_category_id = position_category.id
							where vacancy.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "city_name", "position", "vacancy_description",
							"salary", "employer_id", "work_type_name", "path_to_company_avatar", "created_at", "updated_at", "company_name", "category_name", "compressed_image"}).
							AddRow(1, "Moscow", "Скульптор", "Требуется скульптор без опыта работы", 90000, 1,
								"Полная занятость", "", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "Мэрия Москвы", "Творец", "test"))
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
					mock.ExpectQuery(`select vacancy.id, city.city_name, vacancy.position, vacancy_description, salary, employer_id, work_type.work_type_name, path_to_company_avatar, vacancy.created_at, vacancy.updated_at, 
							company.company_name, position_category.category_name, vacancy.compressed_image from vacancy
							left join work_type on vacancy.work_type_id=work_type.id left join city on vacancy.city_id=city.id
							left join employer on vacancy.employer_id=employer.id left join company on employer.company_name_id=company.id
							left join position_category on vacancy.position_category_id = position_category.id
							where vacancy.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "city_name", "position", "vacancy_description",
							"salary", "employer_id", "work_type_name", "path_to_company_avatar", "created_at", "updated_at", "company_name", "category_name", "compressed_image"}).
							AddRow(1, nil, nil, nil, nil, nil,
								nil, nil, nil, nil, nil, nil, nil))
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

			s := NewVacanciesStorage(db)

			if _, err := s.GetByID(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v", err != nil, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresSearchAll(t *testing.T) {
	t.Parallel()
	type args struct {
		offset                     uint64
		num                        uint64
		searchStr, group, searchBy string
		query                      func(mock sqlmock.Sqlmock, args args)
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
				offset:    0,
				num:       1,
				searchStr: "Скульптор",
				group:     "Творец",
				searchBy:  "position",
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.group,
							args.searchStr,
							args.searchStr,
							args.num,
							args.offset,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "city_name", "position", "vacancy_description",
							"salary", "employer_id", "work_type_name", "path_to_company_avatar", "created_at", "updated_at", "company_name", "category_name", "compressed_image"}).
							AddRow(1, "Moscow", "Скульптор", "Требуется скульптор без опыта работы", 90000, 1,
								"Полная занятость", "", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "Мэрия Москвы", "Творец", "test"))
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

			s := NewVacanciesStorage(db)

			if _, err := s.SearchAll(tt.args.offset, tt.args.num, tt.args.searchStr, tt.args.group, tt.args.searchBy); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v", err != nil, tt.wantErr)
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
		vacancy    dto.JSONVacancy
		worTypeId  uint64
		locationId uint64
		query1     func(mock sqlmock.Sqlmock, args args)
		query2     func(mock sqlmock.Sqlmock, args args)
		query3     func(mock sqlmock.Sqlmock, args args)
		query4     func(mock sqlmock.Sqlmock, args args)
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
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into vacancy (.+)`).
						WithArgs(
							args.vacancy.Position,
							args.vacancy.Description,
							args.vacancy.Salary,
							args.vacancy.EmployerID,
							args.worTypeId,
							args.vacancy.Avatar,
							args.locationId,
							args.vacancy.CompressedAvatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {

				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk1",
			args: args{
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into work_type (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into vacancy (.+)`).
						WithArgs(
							args.vacancy.Position,
							args.vacancy.Description,
							args.vacancy.Salary,
							args.vacancy.EmployerID,
							args.worTypeId,
							args.vacancy.Avatar,
							args.locationId,
							args.vacancy.CompressedAvatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk1",
			args: args{
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into city (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},

				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into vacancy (.+)`).
						WithArgs(
							args.vacancy.Position,
							args.vacancy.Description,
							args.vacancy.Salary,
							args.vacancy.EmployerID,
							args.worTypeId,
							args.vacancy.Avatar,
							args.locationId,
							args.vacancy.CompressedAvatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				vacancy: dto.JSONVacancy{
					EmployerID:           1,
					Salary:               10000,
					Position:             "Position",
					Location:             "Location",
					Description:          "Description",
					WorkType:             "WorkType",
					CompressedAvatar:     "CompressedAvatar",
					PositionCategoryName: "PositionCategoryName",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into vacancy (.+)`).
						WithArgs(
							args.vacancy.Position,
							args.vacancy.Description,
							args.vacancy.Salary,
							args.vacancy.EmployerID,
							args.worTypeId,
							args.vacancy.Avatar,
							args.locationId,
							1,
							args.vacancy.CompressedAvatar,
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
			tt.args.query2(mock, tt.args)
			tt.args.query3(mock, tt.args)
			tt.args.query4(mock, tt.args)

			s := NewVacanciesStorage(db)

			if _, err := s.Create(&tt.args.vacancy); (err != nil) != tt.wantErr {
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
		ID         uint64
		vacancy    dto.JSONVacancy
		worTypeId  uint64
		locationId uint64
		query1     func(mock sqlmock.Sqlmock, args args)
		query2     func(mock sqlmock.Sqlmock, args args)
		query3     func(mock sqlmock.Sqlmock, args args)
		query4     func(mock sqlmock.Sqlmock, args args)
		query5     func(mock sqlmock.Sqlmock, args args)
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
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into work_type  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update vacancy (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
							args.vacancy.Salary,
							args.vacancy.Position,
							args.locationId,
							args.vacancy.Description,
							args.worTypeId,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "position", "vacancy_description",
							"salary", "employer_id", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Скульптор", "Требуется скульптор без опыта работы", 90000, 1, "",
								"2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "test"))
				},
				query5: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"company.company_name"}).
							AddRow("a"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into city (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update vacancy (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
							args.vacancy.Salary,
							args.vacancy.Position,
							args.locationId,
							args.vacancy.Description,
							args.worTypeId,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "position", "vacancy_description",
							"salary", "employer_id", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Скульптор", "Требуется скульптор без опыта работы", 90000, 1, "",
								"2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "test"))
				},
				query5: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"company.company_name"}).
							AddRow("a"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					CompressedAvatar: "CompressedAvatar",
					PositionCategoryName: "PositionCategoryName",
				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update vacancy (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
							args.vacancy.Salary,
							args.vacancy.Position,
							args.locationId,
							args.vacancy.Description,
							args.worTypeId,
							1,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "position", "vacancy_description",
							"salary", "employer_id", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Скульптор", "Требуется скульптор без опыта работы", 90000, 1, "",
								"2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "test"))
				},
				query5: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"company.company_name"}).
							AddRow("a"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					Avatar:           "Avatar",
					CompressedAvatar: "CompressedAvatar",
					PositionCategoryName: "PositionCategoryName",

				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select (.+)`).
						WithArgs(
							args.vacancy.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update vacancy (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
							args.vacancy.Salary,
							args.vacancy.Position,
							args.locationId,
							args.vacancy.Description,
							args.worTypeId,
							args.vacancy.Avatar,
							1,
							args.vacancy.CompressedAvatar,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "position", "vacancy_description",
							"salary", "employer_id", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Скульптор", "Требуется скульптор без опыта работы", 90000, 1, "",
								"2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "test"))
				},
				query5: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"company.company_name"}).
							AddRow("a"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				vacancy: dto.JSONVacancy{
					EmployerID:       1,
					Salary:           10000,
					Position:         "Position",
					Location:         "Location",
					Description:      "Description",
					WorkType:         "WorkType",
					Avatar:           "Avatar",
					CompressedAvatar: "CompressedAvatar",

				},
				worTypeId:  1,
				locationId: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.Location,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.WorkType,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
				},
				query4: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update vacancy (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
							args.vacancy.Salary,
							args.vacancy.Position,
							args.locationId,
							args.vacancy.Description,
							args.worTypeId,
							args.vacancy.Avatar,
							args.vacancy.CompressedAvatar,
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "position", "vacancy_description",
							"salary", "employer_id", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Скульптор", "Требуется скульптор без опыта работы", 90000, 1, "",
								"2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "test"))
				},
				query5: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.vacancy.EmployerID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"company.company_name"}).
							AddRow("a"))
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
			tt.args.query3(mock, tt.args)
			tt.args.query4(mock, tt.args)
			tt.args.query5(mock, tt.args)

			s := NewVacanciesStorage(db)

			if _, err := s.Update(tt.args.ID, &tt.args.vacancy); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresDelete(t *testing.T) {
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
					mock.ExpectExec(`delete from vacancy where id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnResult(sqlmock.NewResult(0, 1))
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

			s := NewVacanciesStorage(db)

			if err := s.Delete(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresSubscribe(t *testing.T) {
	t.Parallel()
	type args struct {
		ID          uint64
		applicantID uint64
		query       func(mock sqlmock.Sqlmock, args args)
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
				ID:          1,
				applicantID: 1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectExec(`insert into vacancy_subscriber (.+)`).
						WithArgs(
							args.ID,
							args.applicantID,
						).
						WillReturnResult(sqlmock.NewResult(1, 0))
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

			s := NewVacanciesStorage(db)

			if err := s.Subscribe(tt.args.ID, tt.args.applicantID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresGetSubscriptionStatus(t *testing.T) {
	t.Parallel()
	type args struct {
		ID          uint64
		applicantID uint64
		query       func(mock sqlmock.Sqlmock, args args)
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
				ID:          1,
				applicantID: 1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select applicant_id from vacancy_subscriber where (.+)`).
						WithArgs(
							args.ID,
							args.applicantID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id"}).
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

			tt.args.query(mock, tt.args)

			s := NewVacanciesStorage(db)

			if _, err := s.GetSubscriptionStatus(tt.args.ID, tt.args.applicantID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v", err != nil, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresGetSubscribersCount(t *testing.T) {
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
					mock.ExpectQuery(`select (.+) from vacancy_subscriber (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).
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

			tt.args.query(mock, tt.args)

			s := NewVacanciesStorage(db)

			if _, err := s.GetSubscribersCount(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresGetSubscribersList(t *testing.T) {
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
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"applicant_id", "first_name", "last_name", "city.city_name",
							"birth_date", "path_to_profile_avatar", "contacts", "education", "email", "password_hash", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, "Иван", "Иванов", "Москва", "12-12-2001", "/src",
								"tg - ", " ", "a@mail.ru", "hash", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "image"))
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

			s := NewVacanciesStorage(db)

			if _, err := s.GetSubscribersList(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresUnsubscribe(t *testing.T) {
	t.Parallel()
	type args struct {
		ID          uint64
		applicantID uint64
		query       func(mock sqlmock.Sqlmock, args args)
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
				ID:          1,
				applicantID: 1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectExec(`delete from vacancy_subscriber where applicant_id=(.+) and vacancy_id=(.+)`).
						WithArgs(
							args.ID,
							args.applicantID,
						).
						WillReturnResult(sqlmock.NewResult(1, 0))
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

			s := NewVacanciesStorage(db)

			if err := s.Unsubscribe(tt.args.ID, tt.args.applicantID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
