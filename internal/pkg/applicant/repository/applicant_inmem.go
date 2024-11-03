package repository

import (
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
)

// DOES NOT SUPPORT ASYNC

// implementation of repository.applicant interface
// in-memory-db
type applicantRepo struct {
	lastID uint64
	data   []*models.Applicant
}

// Initialize new repo
// Returns pointer to it
func NewRepo() *applicantRepo {
	return &applicantRepo{
		lastID: 1,
		data:   make([]*models.Applicant, 0, 10),
	}
}

// Create new applicant into the db
// Accepts pointer to applicant model
// Returns ID of created applicant and error
func (repo *applicantRepo) Create(applicantInput *dto.ApplicantInput) (*models.Applicant, error) {

	applicant := &models.Applicant{
		FirstName:           applicantInput.FirstName,
		LastName:            applicantInput.LastName,
		CityName:            applicantInput.CityName,
		BirthDate:           applicantInput.BirthDate,
		PathToProfileAvatar: applicantInput.PathToProfileAvatar,
		Contacts:            applicantInput.Contacts,
		Education:           applicantInput.Education,
		Email:               applicantInput.Email,
		PasswordHash:        applicantInput.Password,
		CreatedAt:           time.Now().Format(time.RFC3339),
		UpdatedAt:           time.Now().Format(time.RFC3339),
	}
	applicant.ID = repo.lastID
	repo.lastID++
	repo.data = append(repo.data, applicant)
	return applicant, nil
}

// GetByID gets applicant by ID.
// BUG: If user with this ID doesn't exist panic can be raised
func (repo *applicantRepo) GetByID(id uint64) (*models.Applicant, error) {
	if id > repo.lastID {
		// should return error!!
		// remade with own_errors
		return nil, nil
	}
	return repo.data[id], nil
}

func (repo *applicantRepo) GetByEmail(email string) (*models.Applicant, error) {
	for _, applicant := range repo.data {
		if applicant.Email == email {
			return applicant, nil
		}
	}
	return nil, ErrNoUserExist
}
