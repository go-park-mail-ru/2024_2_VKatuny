package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

func TestPostgresGetCVsByApplicantID(t *testing.T) {
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
					mock.ExpectQuery(`select cv.id, applicant_id, position_rus, position_eng, job_search_status_name, cv_description, working_experience,
							path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image  from cv
							left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id
							where cv.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "position_category.category_name", "cv.created_at", "cv.updated_at", "cv.compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Постоянная", "без опыта",
								"cv.svg", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "compressed_image.svg"))
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
					mock.ExpectQuery(`select cv.id, applicant_id, position_rus, position_eng, job_search_status_name, cv_description, working_experience,
							path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image from cv
							left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id
							where cv.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "position_category.category_name", "cv.created_at", "cv.updated_at", "cv.compressed_image"}).
							AddRow(1, nil, nil, nil, nil, nil, nil,
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if _, err := s.GetCVsByApplicantID(context.Background(), tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
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
		cv                dto.JSONCv
		jobSearchStatusID uint64
		query1            func(mock sqlmock.Sqlmock, args args)
		query2            func(mock sqlmock.Sqlmock, args args)
		query3            func(mock sqlmock.Sqlmock, args args)
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
				cv: dto.JSONCv{
					ApplicantID:         1,
					PositionRu:          "Должность",
					PositionEn:          "Position",
					Description:         "Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "Experience",
					Avatar:              "a.svg",
					CompressedAvatar:    "a.svg",
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into job_search_status (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
							args.cv.CompressedAvatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				cv: dto.JSONCv{
					ApplicantID:          1,
					PositionRu:           "Должность",
					PositionEn:           "Position",
					Description:          "Description",
					JobSearchStatusName:  "Ищу",
					WorkingExperience:    "Experience",
					Avatar:               "a.svg",
					CompressedAvatar:     "a.svg",
					PositionCategoryName: "Category",
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
							1,
							args.cv.CompressedAvatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if _, err := s.Create(context.Background(), &tt.args.cv); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!!!!!!!! %s", err != nil, tt.wantErr, err.Error())
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
					mock.ExpectQuery(`select cv.id, applicant_id, position_rus, position_eng, job_search_status.job_search_status_name, cv_description, working_experience,
							path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image from cv
							left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id where cv.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "position_category.category_name", "cv.created_at", "cv.updated_at", "cv.compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Постоянная", "без опыта",
								"cv.svg", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "compressed_image.svg"))
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
					mock.ExpectQuery(`select cv.id, applicant_id, position_rus, position_eng, job_search_status.job_search_status_name, cv_description, working_experience,
							path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image from cv
							left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id where cv.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "position_category.category_name", "cv.created_at", "cv.updated_at", "cv.compressed_image"}).
							AddRow(1, nil, nil, nil, nil, nil, nil,
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if _, err := s.GetByID(context.Background(), tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!! %s", err != nil, tt.wantErr, err)
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
		ID                uint64
		cv                dto.JSONCv
		jobSearchStatusID uint64
		query1            func(mock sqlmock.Sqlmock, args args)
		query2            func(mock sqlmock.Sqlmock, args args)
		query3            func(mock sqlmock.Sqlmock, args args)
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
				cv: dto.JSONCv{
					ApplicantID:         1,
					PositionRu:          "Должность",
					PositionEn:          "Position",
					Description:         "Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "Experience",
					Avatar:              "a.svg",
					CompressedAvatar:    "a.svg",
					ID:                  1,
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into job_search_status  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
							args.cv.CompressedAvatar,
							args.cv.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				cv: dto.JSONCv{
					ApplicantID:          1,
					PositionRu:           "Должность",
					PositionEn:           "Position",
					Description:          "Description",
					JobSearchStatusName:  "Ищу",
					WorkingExperience:    "Experience",
					Avatar:               "a.svg",
					CompressedAvatar:     "a.svg",
					PositionCategoryName: "Category",
					ID:                   1,
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select id from position_category (.+)`).
						WithArgs(
							args.cv.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
							1,
							args.cv.CompressedAvatar,
							args.cv.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				cv: dto.JSONCv{
					ApplicantID:         1,
					PositionRu:          "Должность",
					PositionEn:          "Position",
					Description:         "Description",
					JobSearchStatusName: "Ищу",
					WorkingExperience:   "Experience",
					CompressedAvatar:    "a.svg",
					ID:                  1,
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnError(sql.ErrNoRows)
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`insert into job_search_status  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "TestOk",
			args: args{
				ID: 1,
				cv: dto.JSONCv{
					ApplicantID:          1,
					PositionRu:           "Должность",
					PositionEn:           "Position",
					Description:          "Description",
					JobSearchStatusName:  "Ищу",
					WorkingExperience:    "Experience",
					CompressedAvatar:     "a.svg",
					PositionCategoryName: "Category",
					ID:                   1,
				},
				jobSearchStatusID: 1,
				query1: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select  (.+)`).
						WithArgs(
							args.cv.JobSearchStatusName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query2: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select id from position_category (.+)`).
						WithArgs(
							args.cv.PositionCategoryName,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(1))
				},
				query3: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`update cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							1,
							args.cv.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at", "compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "cv.svg"))
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if _, err := s.Update(context.Background(), tt.args.ID, &tt.args.cv); (err != nil) != tt.wantErr {
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
					mock.ExpectExec(`delete from cv where id = (.+)`).
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if err := s.Delete(context.Background(), tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
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
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "position_category.category_name", "cv.created_at", "cv.updated_at", "cv.compressed_image"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Постоянная", "без опыта",
								"cv.svg", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300", "compressed_image.svg"))
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

			s := NewCVStorage(db, logrus.StandardLogger())

			if _, err := s.SearchAll(context.Background(), tt.args.offset, tt.args.num, tt.args.searchStr, tt.args.group, tt.args.searchBy); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err !!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
