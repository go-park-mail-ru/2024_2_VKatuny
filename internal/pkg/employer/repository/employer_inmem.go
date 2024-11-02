package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// DOES NOT SUPPORT ASYNC

// implementation of repository.Employer interface
// temp in-memory-db
type employerRepo struct {
	lastID uint64
	data   []*models.Employer
}

// DOES NOT SUPPORT ASYNC
// Initialize new repo
// Returns pointer to it
func NewRepo() *employerRepo {
	return &employerRepo{
		lastID: 0,
		data:   make([]*models.Employer, 0, 10),
	}
}

// Creates new employer
// Accepts pointer to employer model
// Returns ID of created employer and error
func (repo *employerRepo) Create(employerInput *dto.JSONEmployerRegistrationForm) (uint64, error) {
	employer := &models.Employer{
		ID:          repo.lastID,
		Name:        employerInput.Name,
		LastName:    employerInput.LastName,
		Position:    employerInput.Position,
		CompanyName: employerInput.CompanyName,
		Description: employerInput.Description,
		Website:     employerInput.Website,
		Email:       employerInput.Email,
		Password:    employerInput.Password,
	}
	repo.lastID++
	repo.data = append(repo.data, employer)
	return employer.ID, nil
}

// GetByID gets employer by ID
// Accepts ID
// Returns pointer to employer and error
// !!!if NOT FOUND DOESN'T return error yet!!!
func (repo *employerRepo) GetByID(id uint64) (*models.Employer, error) {
	if id > repo.lastID {
		return nil, nil
	}
	return repo.data[id], nil
}

func (repo *employerRepo) GetByEmail(email string) (*models.Employer, error) {
	for _, employer := range repo.data {
		if employer.Email == email {
			return employer, nil
		}
	}
	return nil, employer.ErrNoUserExist
}
