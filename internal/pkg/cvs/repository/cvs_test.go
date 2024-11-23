package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
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
						path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id 
						where cv.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "cv.created_at", "cv.updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Постоянная", "без опыта",
								"cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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
						path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id 
						where cv.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"cv.id", "applicant_id", "position_rus", "position_eng", "cv_description",
							"job_search_status_name", "working_experience", "path_to_profile_avatar", "cv.created_at", "cv.updated_at"}).
							AddRow(1, nil, nil, nil, nil, nil, nil,
								nil, nil, nil))
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

			s := NewCVStorage(db)

			if _, err := s.GetCVsByApplicantID(tt.args.ID); (err != nil) != tt.wantErr {
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
					mock.ExpectQuery(`insert into cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewCVStorage(db)

			if _, err := s.Create(&tt.args.cv); (err != nil) != tt.wantErr {
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
						path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id where cv.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng", "job_search_status.job_search_status_name",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Ищу", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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
						path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id where cv.id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng", "job_search_status.job_search_status_name",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, nil, nil, nil, nil, nil, nil, nil, nil, nil))
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

			s := NewCVStorage(db)

			if _, err := s.GetByID(tt.args.ID); (err != nil) != tt.wantErr {
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
					ID:                  1,
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
					mock.ExpectQuery(`update cv (.+)`).
						WithArgs(
							args.cv.ApplicantID,
							args.cv.PositionRu,
							args.cv.PositionEn,
							args.cv.Description,
							args.jobSearchStatusID,
							args.cv.WorkingExperience,
							args.cv.Avatar,
							args.cv.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewCVStorage(db)

			if _, err := s.Update(tt.args.ID, &tt.args.cv); (err != nil) != tt.wantErr {
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

			s := NewCVStorage(db)

			if err := s.Delete(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!!! %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresGetWithOffset(t *testing.T) {
	t.Parallel()
	type args struct {
		offset uint64
		num    uint64
		query  func(mock sqlmock.Sqlmock, args args)
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
				offset: 0,
				num:    1,
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select cv.id, applicant_id, cv.position_rus, cv.position_eng, cv_description, job_search_status.job_search_status_name,
						working_experience, path_to_profile_avatar, cv.created_at, cv.updated_at from cv
						left join job_search_status on cv.job_search_status_id=job_search_status.id
						ORDER BY created_at desc limit (.+) offset (.+)`).
						WithArgs(
							args.offset,
							args.num,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "job_search_status_name", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Ищу работу", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewCVStorage(db)

			if _, err := s.GetWithOffset(tt.args.num, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresSearchByPositionDescription(t *testing.T) {
	t.Parallel()
	type args struct {
		offset uint64
		num    uint64
		search string
		query  func(mock sqlmock.Sqlmock, args args)
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
				offset: 0,
				num:    1,
				search: "Скульптор",
				query: func(mock sqlmock.Sqlmock, args args) {
					mock.ExpectQuery(`select cv.id, applicant_id, cv.position_rus, cv.position_eng, cv_description, job_search_status.job_search_status_name,
						working_experience, path_to_profile_avatar, cv.created_at, cv.updated_at from cv
						left join job_search_status on cv.job_search_status_id=job_search_status.id
						where (.+)`).
						WithArgs(
							args.search,
							args.search,
							args.num,
							args.offset,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "position_rus", "position_eng",
							"cv_description", "job_search_status_name", "working_experience", "path_to_profile_avatar", "created_at", "updated_at"}).
							AddRow(1, 1, "Скульптор", "Sculptor", "Я усердный и трудолюбивый", "Ищу работу", "Нет опыта", "cv.svg", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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

			s := NewCVStorage(db)

			if _, err := s.SearchByPositionDescription(tt.args.offset, tt.args.num, tt.args.search); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err %s", err != nil, tt.wantErr, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
