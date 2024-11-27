package repository

import (
	"database/sql"
	"fmt"
	"strconv"

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
		path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image  from cv
		left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id
		where cv.applicant_id = $1`, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneCV dto.JSONCvWithNull
		if err := rows.Scan(&oneCV.ID,
			&oneCV.ApplicantID,
			&oneCV.PositionRu,
			&oneCV.PositionEn,
			&oneCV.JobSearchStatusName,
			&oneCV.Description,
			&oneCV.WorkingExperience,
			&oneCV.Avatar,
			&oneCV.PositionCategoryName,
			&oneCV.CreatedAt,
			&oneCV.UpdatedAt,
			&oneCV.CompressedAvatar,
			); err != nil {
			return nil, err
		}
		oneCVOk := dto.JSONCv{
			ID:                   oneCV.ID,
			ApplicantID:          oneCV.ApplicantID,
			PositionRu:           oneCV.PositionRu,
			PositionEn:           oneCV.PositionEn,
			JobSearchStatusName:  oneCV.JobSearchStatusName,
			Description:          oneCV.Description,
			WorkingExperience:    oneCV.WorkingExperience,
			Avatar:               oneCV.Avatar,
			PositionCategoryName: oneCV.PositionCategoryName.String,
			CreatedAt:            oneCV.CreatedAt,
			UpdatedAt:            oneCV.UpdatedAt,
			CompressedAvatar:     oneCV.CompressedAvatar.String,
		}
		CVs = append(CVs, &oneCVOk)
		fmt.Println(oneCVOk)
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
	if cv.PositionCategoryName == "" {
		row = s.db.QueryRow(`insert into cv (applicant_id, position_rus, position_eng, cv_description, job_search_status_id, working_experience, path_to_profile_avatar, compressed_image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id, applicant_id, position_rus, position_eng,
		cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
			cv.ApplicantID, cv.PositionRu, cv.PositionEn, cv.Description, JobSearchStatusID, cv.WorkingExperience, cv.Avatar, cv.CompressedAvatar)
	} else {
		var PositionCategoryID int
		row = s.db.QueryRow(`select id from position_category where category_name=$1`, cv.PositionCategoryName)
		err := row.Scan(&PositionCategoryID)
		if err != nil {
			return nil, err
		}
		row = s.db.QueryRow(`insert into cv (applicant_id, position_rus, position_eng, cv_description, job_search_status_id, working_experience, path_to_profile_avatar, position_category_id, compressed_image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id, applicant_id, position_rus, position_eng,
		cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
			cv.ApplicantID, cv.PositionRu, cv.PositionEn, cv.Description, JobSearchStatusID, cv.WorkingExperience, cv.Avatar, PositionCategoryID, cv.CompressedAvatar)
	}
	err := row.Scan(&oneCv.ID,
		&oneCv.ApplicantID,
		&oneCv.PositionRu,
		&oneCv.PositionEn,
		&oneCv.Description,
		&oneCv.WorkingExperience,
		&oneCv.Avatar,
		&oneCv.CreatedAt,
		&oneCv.UpdatedAt,
		&oneCv.CompressedAvatar,
	)
	oneCv.JobSearchStatusName = cv.JobSearchStatusName
	oneCv.PositionCategoryName = cv.PositionCategoryName
	if err != nil {
		return nil, err
	}
	return &oneCv, err
}

func (s *PostgreSQLCVStorage) GetByID(ID uint64) (*dto.JSONCv, error) {

	row := s.db.QueryRow(`select cv.id, applicant_id, position_rus, position_eng, job_search_status.job_search_status_name, cv_description, working_experience,
		path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image from cv
		left join job_search_status on job_search_status.id = cv.job_search_status_id left join position_category on cv.position_category_id = position_category.id where cv.id = $1`, ID)
	var oneCV dto.JSONCvWithNull

	err := row.Scan(
		&oneCV.ID,
		&oneCV.ApplicantID,
		&oneCV.PositionRu,
		&oneCV.PositionEn,
		&oneCV.JobSearchStatusName,
		&oneCV.Description,
		&oneCV.WorkingExperience,
		&oneCV.Avatar,
		&oneCV.PositionCategoryName,
		&oneCV.CreatedAt,
		&oneCV.UpdatedAt,
		&oneCV.CompressedAvatar,
	)

	oneCVOk := dto.JSONCv{
		ID:                   oneCV.ID,
		ApplicantID:          oneCV.ApplicantID,
		PositionRu:           oneCV.PositionRu,
		PositionEn:           oneCV.PositionEn,
		Description:          oneCV.Description,
		JobSearchStatusName:  oneCV.JobSearchStatusName,
		WorkingExperience:    oneCV.WorkingExperience,
		Avatar:               oneCV.Avatar,
		PositionCategoryName: oneCV.PositionCategoryName.String,
		CreatedAt:            oneCV.CreatedAt,
		UpdatedAt:            oneCV.UpdatedAt,
		CompressedAvatar:     oneCV.CompressedAvatar.String,
	}
	if err != nil {
		return nil, err
	}
	return &oneCVOk, err
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
	var PositionCategoryID int
	if updatedCv.PositionCategoryName != "" {
		row = s.db.QueryRow(`select id from position_category where category_name=$1`, updatedCv.PositionCategoryName)
		err := row.Scan(&PositionCategoryID)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(updatedCv)
	if updatedCv.Avatar != "" {
		if updatedCv.PositionCategoryName == "" {
			row = s.db.QueryRow(`update cv
				set applicant_id = $1, position_rus = $2, position_eng = $3, cv_description=$4, 
				job_search_status_id = $5, working_experience = $6, path_to_profile_avatar=$7, compressed_image=$8 where id=$9 returning id, 
				applicant_id, position_rus, position_eng, cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
				updatedCv.ApplicantID, updatedCv.PositionRu, updatedCv.PositionEn, updatedCv.Description, JobSearchStatusID, updatedCv.WorkingExperience, updatedCv.Avatar, updatedCv.CompressedAvatar, ID)
		} else {
			row = s.db.QueryRow(`update cv
				set applicant_id = $1, position_rus = $2, position_eng = $3, cv_description=$4, 
				job_search_status_id = $5, working_experience = $6, path_to_profile_avatar=$7, position_category_id=$8, compressed_image=$9 where id=$10 returning id, 
				applicant_id, position_rus, position_eng, cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
				updatedCv.ApplicantID, updatedCv.PositionRu, updatedCv.PositionEn, updatedCv.Description, JobSearchStatusID, updatedCv.WorkingExperience, updatedCv.Avatar, PositionCategoryID, updatedCv.CompressedAvatar, ID)
		}
	} else {
		if updatedCv.PositionCategoryName == "" {
			row = s.db.QueryRow(`update cv
				set applicant_id = $1, position_rus = $2, position_eng = $3, cv_description=$4, 
				job_search_status_id = $5, working_experience = $6 where id=$7 returning id, 
				applicant_id, position_rus, position_eng, cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
				updatedCv.ApplicantID, updatedCv.PositionRu, updatedCv.PositionEn, updatedCv.Description, JobSearchStatusID, updatedCv.WorkingExperience, ID)
		} else {
			row = s.db.QueryRow(`update cv
					set applicant_id = $1, position_rus = $2, position_eng = $3, cv_description=$4, 
					job_search_status_id = $5, working_experience = $6, path_to_profile_avatar=$7, position_category_id=$8 where id=$9 returning id, 
					applicant_id, position_rus, position_eng, cv_description, working_experience, path_to_profile_avatar, created_at, updated_at, compressed_image`,
				updatedCv.ApplicantID, updatedCv.PositionRu, updatedCv.PositionEn, updatedCv.Description, JobSearchStatusID, updatedCv.WorkingExperience, updatedCv.Avatar, PositionCategoryID, ID)
		}
	}

	var oneCV dto.JSONCvWithNull

	err := row.Scan(
		&oneCV.ID,
		&oneCV.ApplicantID,
		&oneCV.PositionRu,
		&oneCV.PositionEn,
		&oneCV.Description,
		&oneCV.WorkingExperience,
		&oneCV.Avatar,
		&oneCV.CreatedAt,
		&oneCV.UpdatedAt,
		&oneCV.CompressedAvatar,
	)

	oneCVOk := dto.JSONCv{
		ID:                   oneCV.ID,
		ApplicantID:          oneCV.ApplicantID,
		PositionRu:           oneCV.PositionRu,
		PositionEn:           oneCV.PositionEn,
		Description:          oneCV.Description,
		WorkingExperience:    oneCV.WorkingExperience,
		Avatar:               oneCV.Avatar,
		CreatedAt:            oneCV.CreatedAt,
		UpdatedAt:            oneCV.UpdatedAt,
		CompressedAvatar:     oneCV.CompressedAvatar.String,
	}
	oneCVOk.JobSearchStatusName = updatedCv.JobSearchStatusName
	oneCVOk.PositionCategoryName = updatedCv.PositionCategoryName
	if err != nil {
		return nil, err
	}
	return &oneCVOk, err
}

func (s *PostgreSQLCVStorage) Delete(ID uint64) error {
	_, err := s.db.Exec(`delete from cv where id = $1`, ID)
	return err
}

func (s *PostgreSQLCVStorage) SearchAll(offset uint64, num uint64, searchStr, group, searchBy string) ([]*dto.JSONCv, error) {
	CVs := make([]*dto.JSONCv, 0)
	iter := 1
	mainPart := `select cv.id, applicant_id, cv.position_rus, cv.position_eng, cv_description, job_search_status.job_search_status_name,
		working_experience, path_to_profile_avatar, position_category.category_name, cv.created_at, cv.updated_at, cv.compressed_image from cv
		left join job_search_status on cv.job_search_status_id=job_search_status.id
		left join position_category on cv.position_category_id = position_category.id `
	categoryPart := ""
	if group != "" {
		categoryPart = "where position_category.category_name = $" + strconv.Itoa(iter)
		iter++
	}
	searchPart := ""
	if searchStr != "" {
		if iter != 1 {
			searchPart += " and "
		} else {
			searchPart += " where "
		}
		weights := "'{1, 1, 1, 1}'"
		language := "'russian'"
		switch searchBy {
		case "position_rus":
			weights = "'{0, 0, 0, 1}'"
		case "position_eng":
			weights = "'{0, 0, 1, 0}'"
			language = "'english'"
		case "working_experience":
			weights = "'{0, 1, 0, 0}'"
		case "cv_description":
			weights = "'{1, 0, 0, 0}'"
		}
		searchPart += "ts_rank_cd(" + weights + ", cv.fts, plainto_tsquery(" + language + ", $" + strconv.Itoa(iter) + ")) <> 0 order by ts_rank_cd(" + weights + ", cv.fts, plainto_tsquery(" + language + ", $" + strconv.Itoa(iter+1) + ")) desc "
		iter += 2
	} else {
		searchPart += " ORDER BY created_at desc "
	}

	lastPart := " limit $" + strconv.Itoa(iter) + " offset $" + strconv.Itoa(iter+1)
	fmt.Println(categoryPart)
	fmt.Println(mainPart + categoryPart + searchPart + lastPart)
	iter += 2
	var rows *sql.Rows
	var err error
	if group != "" && searchStr != "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, group, searchStr, searchStr, num, offset)
	} else if group == "" && searchStr != "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, searchStr, searchStr, num, offset)
	} else if group == "" && searchStr == "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, num, offset)
	} else if group != "" && searchStr == "" {
		rows, err = s.db.Query(mainPart+categoryPart+searchPart+lastPart, group, num, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var oneCV dto.JSONCvWithNull
		if err := rows.Scan(
			&oneCV.ID,
			&oneCV.ApplicantID,
			&oneCV.PositionRu,
			&oneCV.PositionEn,
			&oneCV.JobSearchStatusName,
			&oneCV.Description,
			&oneCV.WorkingExperience,
			&oneCV.Avatar,
			&oneCV.PositionCategoryName,
			&oneCV.CreatedAt,
			&oneCV.UpdatedAt,
			&oneCV.CompressedAvatar,
			); err != nil {
			return nil, err
		}
		oneCVOk := dto.JSONCv{
			ID:                   oneCV.ID,
			ApplicantID:          oneCV.ApplicantID,
			PositionRu:           oneCV.PositionRu,
			PositionEn:           oneCV.PositionEn,
			JobSearchStatusName:  oneCV.JobSearchStatusName,
			Description:          oneCV.Description,
			WorkingExperience:    oneCV.WorkingExperience,
			Avatar:               oneCV.Avatar,
			PositionCategoryName: oneCV.PositionCategoryName.String,
			CreatedAt:            oneCV.CreatedAt,
			UpdatedAt:            oneCV.UpdatedAt,
			CompressedAvatar:     oneCV.CompressedAvatar.String,
		}
		CVs = append(CVs, &oneCVOk)
		fmt.Println(oneCVOk)
	}
	fmt.Println(CVs)
	return CVs, nil
}
