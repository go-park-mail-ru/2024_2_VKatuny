package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
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
	return &vacanciesRepo{
		lastID: 1, // for oleg's db
		data: make([]*models.Vacancy, 0, 10),
	}
}

// Add new vacncy into the db
// Accepts pointer to vacancy model
// Returns ID of created vacancy and error
func (repo *vacanciesRepo) Add(vacancy *models.Vacancy) (uint64, error) {
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
