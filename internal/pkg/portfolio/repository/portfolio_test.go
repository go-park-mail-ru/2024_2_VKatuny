package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
					mock.ExpectQuery(`select id, applicant_id, portfolio_name, created_at, updated_at  from portfolio where portfolio.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "portfolio_name", "created_at", "updated_at"}).
							AddRow(1, 1, "Portfolio", "2024-11-09 04:17:52.598 +0300", "2024-11-09 04:17:52.598 +0300"))
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
					mock.ExpectQuery(`select id, applicant_id, portfolio_name, created_at, updated_at  from portfolio where portfolio.applicant_id = (.+)`).
						WithArgs(
							args.ID,
						).
						WillReturnRows(sqlmock.NewRows([]string{"id", "applicant_id", "portfolio_name", "created_at", "updated_at"}).
							AddRow(1, nil, nil, nil, nil))
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

			s := NewPortfolioStorage(db)

			if _, err := s.GetPortfoliosByApplicantID(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Postgres error = %v, wantErr %v, err!!! %s", err != nil, tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
