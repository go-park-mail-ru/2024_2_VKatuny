package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// DOES NOT SUPPORT ASYNC

// implementation of repository.Vacancies interface
// in-memory-db
type vacanciesRepo struct {
	lastID uint64
	data   []*models.Vacancy
}

// Initialize new repo
// Returns pointer to it
func NewRepo() *vacanciesRepo {
	vacancies := &vacanciesRepo{
		lastID: 0,
		data:   make([]*models.Vacancy, 0, 10),
	}
	vacancies.lastID = 25
	for i := uint64(0); i < 25; i += 5 {
		vacancies.data = append(vacancies.data, &models.Vacancy{
			ID:       i,
			Position: "Продавец консультант",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary:     "Не указана",
			EmployerID: 1,
			CreatedAt:  "2024.09.29 16:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:       "img/picture_name1.png",
		})
		vacancies.data = append(vacancies.data, &models.Vacancy{
			ID:       i + 1,
			Position: "Продавец",
			Description: `Ищем продавца на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую работу. Своевременную оплату гарантируем.`,
			Salary:     "80 000",
			EmployerID: 1,
			CreatedAt:  "2024.09.29 17:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:       "img/picture_name2.png",
		})
		vacancies.data = append(vacancies.data, &models.Vacancy{
			ID:       i + 2,
			Position: "Администратор",
			Description: `Ищем администратора на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на продуктивную работу с людьми. Своевременную оплату гарантируем.`,
			Salary:     "100 500",
			EmployerID: 1,
			CreatedAt:  "2024.09.29 18:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:       "img/picture_name3.png",
		})
		vacancies.data = append(vacancies.data, &models.Vacancy{
			ID:       i + 3,
			Position: "Охранник",
			Description: `Ищем охранника на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую посменную работу. Своевременную оплату гарантируем.`,
			Salary:     "Не указана",
			EmployerID: 1,
			CreatedAt:  "2024.09.29 19:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:       "img/picture_name4.png",
		})
		vacancies.data = append(vacancies.data, &models.Vacancy{
			ID:       i + 4,
			Position: "Уборщик помещений",
			Description: `Ищем уборщика на полную ставку в ближайший магазин.
			Требуются ответственные личности, способные на тяжелую кропотливую работу. Своевременную оплату гарантируем.`,
			Salary:     "50 000",
			EmployerID: 1,
			CreatedAt:  "2024.09.29 20:55:00", // YYYY.MM.DD HH:MM:SS
			Logo:       "img/picture_name5.png",
		})
	}
	return vacancies
}

// Add new vacncy into the db
// Accepts pointer to vacancy model
// Returns ID of created vacancy and error
func (repo *vacanciesRepo) Create(vacancy *models.Vacancy) (uint64, error) {
	repo.lastID++
	vacancy.ID = repo.lastID
	repo.data = append(repo.data, vacancy)
	return vacancy.ID, nil
}

// GetSomeVacancies get some num amount of vacancies from db starting from offset.
// Dosn't support case when there aren't at least one element in range [offset, offset + num).
// In this case method way cause PANIC.
func (repo *vacanciesRepo) GetWithOffset(offset uint64, num uint64) ([]*models.Vacancy, error) {
	leftBound := offset
	rightBound := offset + num
	// covering cases when offset is out of slice bounds
	if leftBound > repo.lastID {
		rightBound = leftBound
	} else if rightBound > repo.lastID {
		rightBound = repo.lastID
	}
	vacancies := repo.data[leftBound:rightBound]
	return vacancies, nil
}
