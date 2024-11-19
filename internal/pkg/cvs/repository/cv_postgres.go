package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

type PostgreSQLCVStorage struct {
	db *sql.DB
}

func NewCVStorage(db *sql.DB) *PostgreSQLCVStorage {
	return &PostgreSQLCVStorage{
		db: db,
	}
}

func (s *PostgreSQLCVStorage) GetCVsByApplicantID(applicantID uint64) ([]*dto.JSONCv, error) {

	CVs := make([]*dto.JSONCv, 0)

	rows, err := s.db.Query(`select cv.id, applicant_id, position_rus, position_eng, job_search_status_name, cv_description, working_experience,
		path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id 
		where cv.applicant_id = $1`, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var CV dto.JSONCv
		if err := rows.Scan(&CV.ID, &CV.ApplicantID, &CV.PositionRu, &CV.PositionEn, &CV.JobSearchStatusName, &CV.Description, &CV.WorkingExperience, &CV.Avatar, &CV.CreatedAt, &CV.UpdatedAt); err != nil {
			return nil, err
		}
		CVs = append(CVs, &CV)
		fmt.Println(CV)
	}

	return CVs, nil
}

func (s *PostgreSQLCVStorage) Create(cv *dto.JSONCv) (*dto.JSONCv, error) {
	var JobSearchStatusID int
	row := s.db.QueryRow(`select id from job_search_status where job_search_status_name=$1`, cv.JobSearchStatusName)
	if err := row.Scan(&JobSearchStatusID); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into job_search_status (job_search_status_name) VALUES ($1) returning id`, cv.JobSearchStatusName)
			err = row.Scan(&JobSearchStatusID)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	var oneCv dto.JSONCv
	row = s.db.QueryRow(`insert into cv (applicant_id, position_rus, position_eng, cv_description, job_search_status_id, working_experience)
		VALUES ($1, $2, $3, $4, $5, $6) returning id, applicant_id, position_rus, position_eng,
		cv_description, working_experience, path_to_profile_avatar, created_at, updated_at`,
		cv.ApplicantID, cv.PositionRu, cv.PositionEn, cv.Description, JobSearchStatusID, cv.WorkingExperience)
	err := row.Scan(&oneCv.ID,
		&oneCv.ApplicantID,
		&oneCv.PositionRu,
		&oneCv.PositionEn,
		&oneCv.Description,
		&oneCv.WorkingExperience,
		&oneCv.Avatar,
		&oneCv.CreatedAt,
		&oneCv.UpdatedAt)
	oneCv.JobSearchStatusName = cv.JobSearchStatusName
	if err != nil {
		return nil, err
	}
	return &oneCv, err
}

func (s *PostgreSQLCVStorage) GetByID(ID uint64) (*dto.JSONCv, error) {

	row := s.db.QueryRow(`select cv.id, applicant_id, position_rus, position_eng, job_search_status.job_search_status_name, cv_description, working_experience,
		path_to_profile_avatar, cv.created_at, cv.updated_at  from cv left join job_search_status on job_search_status.id = cv.job_search_status_id where cv.id = $1`, ID)
	var oneCv dto.JSONCv

	err := row.Scan(&oneCv.ID,
		&oneCv.ApplicantID,
		&oneCv.PositionRu,
		&oneCv.PositionEn,
		&oneCv.JobSearchStatusName,
		&oneCv.Description,
		&oneCv.WorkingExperience,
		&oneCv.Avatar,
		&oneCv.CreatedAt,
		&oneCv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &oneCv, err
}

func (s *PostgreSQLCVStorage) Update(ID uint64, updatedCv *dto.JSONCv) (*dto.JSONCv, error) {

	var JobSearchStatusID int
	row := s.db.QueryRow(`select id from job_search_status where job_search_status_name=$1`, updatedCv.JobSearchStatusName)
	if err := row.Scan(&JobSearchStatusID); err != nil {
		switch err {
		case sql.ErrNoRows:
			row = s.db.QueryRow(`insert into job_search_status (job_search_status_name) VALUES ($1) returning id`, updatedCv.JobSearchStatusName)
			err = row.Scan(&JobSearchStatusID)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	row = s.db.QueryRow(`update cv
		set applicant_id = $1, position_rus = $2, position_eng = $3, cv_description=$4, 
		job_search_status_id = $5, working_experience = $6, path_to_profile_avatar=$7 where id=$8 returning id, 
		applicant_id, position_rus, position_eng, cv_description, working_experience, path_to_profile_avatar, created_at, updated_at`,
		updatedCv.ApplicantID, updatedCv.PositionRu, updatedCv.PositionEn, updatedCv.Description, JobSearchStatusID, updatedCv.WorkingExperience, updatedCv.Avatar, ID)

	var oneCV dto.JSONCv

	err := row.Scan(&oneCV.ID,
		&oneCV.ApplicantID,
		&oneCV.PositionRu,
		&oneCV.PositionEn,
		&oneCV.Description,
		&oneCV.WorkingExperience,
		&oneCV.Avatar,
		&oneCV.CreatedAt,
		&oneCV.UpdatedAt)
	oneCV.JobSearchStatusName = updatedCv.JobSearchStatusName
	if err != nil {
		return nil, err
	}
	return &oneCV, err
}

func (s *PostgreSQLCVStorage) Delete(ID uint64) error {
	_, err := s.db.Exec(`delete from cv where id = $1`, ID)
	return err
}

func (s *PostgreSQLCVStorage) GetWithOffset(offset uint64, num uint64) ([]*dto.JSONCv, error) {
	CVs := make([]*dto.JSONCv, 0)
	rows, err := s.db.Query(`select cv.id, applicant_id, cv.position_rus, cv.position_eng, cv_description, job_search_status.job_search_status_name,
		working_experience, path_to_profile_avatar, cv.created_at, cv.updated_at from cv
		left join job_search_status on cv.job_search_status_id=job_search_status.id
		ORDER BY created_at desc limit $1 offset $2`, num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var CV dto.JSONCv
		if err := rows.Scan(&CV.ID, &CV.ApplicantID, &CV.PositionRu, &CV.PositionEn, &CV.Description, &CV.JobSearchStatusName, &CV.WorkingExperience,
			&CV.Avatar, &CV.CreatedAt, &CV.UpdatedAt); err != nil {
			return nil, err
		}
		CVs = append(CVs, &CV)
		fmt.Println(CV)
	}

	return CVs, nil
}

func (s *PostgreSQLCVStorage) SearchByPositionDescription(offset uint64, num uint64, searchStr string) ([]*dto.JSONCv, error) {
	CVs := make([]*dto.JSONCv, 0)
	rows, err := s.db.Query(`select cv.id, applicant_id, cv.position_rus, cv.position_eng, cv_description, job_search_status.job_search_status_name,
		working_experience, path_to_profile_avatar, cv.created_at, cv.updated_at from cv
		left join job_search_status on cv.job_search_status_id=job_search_status.id
		where ts_rank_cd(fts, plainto_tsquery('russian', $1)) <> 0 order by ts_rank_cd(fts, plainto_tsquery('russian', $2)) desc limit $3 offset $4`, searchStr, searchStr, num, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var CV dto.JSONCv
		if err := rows.Scan(&CV.ID, &CV.ApplicantID, &CV.PositionRu, &CV.PositionEn, &CV.Description, &CV.JobSearchStatusName, &CV.WorkingExperience,
			&CV.Avatar, &CV.CreatedAt, &CV.UpdatedAt); err != nil {
			return nil, err
		}
		CVs = append(CVs, &CV)
		fmt.Println(CV)
	}
	fmt.Println(CVs)
	return CVs, nil
}
