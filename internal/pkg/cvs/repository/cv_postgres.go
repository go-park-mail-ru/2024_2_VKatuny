package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

type PostgreSQLCVStorage struct {
	db *sql.DB
}

func NewCVStorage(db *sql.DB) *PostgreSQLCVStorage {
	return &PostgreSQLCVStorage{
		db: db,
	}
}

func (s *PostgreSQLCVStorage) GetCVsByApplicantID(applicantID uint64) ([]*models.CV, error) {

	CVs := make([]*models.CV, 0)

	rows, err := s.db.Query(`select id, applicant_id, position_rus, position_eng, job_search_status_id, working_experience,
		path_to_profile_avatar, created_at, updated_at  from cv where cv.applicant_id = $1`, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var CV models.CV
		if err := rows.Scan(&CV.ID, &CV.ApplicantID, &CV.PositionRus, &CV.PositionEng, &CV.JobSearchStatusID, &CV.WorkingExperience, &CV.PathToProfileAvatar, &CV.CreatedAt, &CV.UpdatedAt); err != nil {
			return nil, err
		}
		fmt.Println(CV)
	}

	return CVs, nil
}
