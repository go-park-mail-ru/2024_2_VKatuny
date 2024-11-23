package repository

import (
	"database/sql"
	"fmt"

	dto "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/survey/survey"
)

type PostgreSQLSurveyRepository struct {
	db *sql.DB
}

func NewPostgreSQLSurveyRepository(db *sql.DB) *PostgreSQLSurveyRepository {
	return &PostgreSQLSurveyRepository{
		db: db,
	}
}

func (sr *PostgreSQLSurveyRepository) GetStatistic() ([]*dto.Statistics, error) {
	var statisticsOutput []*dto.Statistics
	rows, err := sr.db.Query(`select AVG(val), question.question_text, question.id from answer
	left join question on question.id = answer.question_id
	left join question_type on question_type.id = question.type_id
	group by question.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stattistic dto.Statistics
		if err := rows.Scan(&stattistic.ValAVG, &stattistic.QuestionText, &stattistic.QuestionID); err != nil {
			return nil, err
		}
		statisticsOutput = append(statisticsOutput, &stattistic)
	}

	return statisticsOutput, nil
}

func (sr *PostgreSQLSurveyRepository) GetQuestionByType() ([]*dto.Question, error) {
	Questions := make([]*dto.Question, 0)
	rows, err := sr.db.Query(`select question.id, question.question_text, question_type.question_type_name, question.position from question
	left join question_type on question_type.id = question.type_id
	order by position ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Question dto.Question
		if err := rows.Scan(&Question.ID, &Question.QuestionText, &Question.TypeText, &Question.Position); err != nil {
			return nil, err
		}
		Questions = append(Questions, &Question)
		fmt.Println(Question)
	}

	return Questions, nil
}

func (sr *PostgreSQLSurveyRepository) CreateAnswerAuthorised(QuestionAnswer *dto.QuestionAnswer) error {
	var UserID int
	row := sr.db.QueryRow(`select id from user where token=$1`, QuestionAnswer.Token)
	err := row.Scan(&UserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			row = sr.db.QueryRow(`insert into user (token) VALUES ($1) returning id`, QuestionAnswer.Token)
			err = row.Scan(&UserID)
			if err != nil {
				return err
			}
		default:
			return err
		}
	} else {
		return nil // already vouted
	}
	row = sr.db.QueryRow(`insert into answer (user_id, value, question_id) VALUES ($1, $2, $3)`,
		UserID, QuestionAnswer.Value, QuestionAnswer.QuestionID)
	return nil
}

// func (sr *PostgreSQLSurveyRepository) CreateAnswerUnauthorised(QuestionAnswer *dto.QuestionAnswer) error {
// 	row := s.db.QueryRow(`insert into answer (value, question_id) VALUES ($1, $2)`,
// 	QuestionAnswer.Value, QuestionAnswer.QuestionID)
// 	return nil
// }
