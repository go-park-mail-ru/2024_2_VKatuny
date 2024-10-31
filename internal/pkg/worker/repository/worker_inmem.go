package repository

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker"
)

// DOES NOT SUPPORT ASYNC

// implementation of repository.Worker interface
// in-memory-db
type workerRepo struct {
	lastID uint64
	data   []*models.Worker
}

// Initialize new repo
// Returns pointer to it
func NewRepo() *workerRepo {
	return &workerRepo{
		lastID: 1,
		data:   make([]*models.Worker, 0, 10),
	}
}

// Add new worker into the db
// Accepts pointer to worker model
// Returns ID of created worker and error
func (repo *workerRepo) Add(worker *models.Worker) (uint64, error) {
	worker.ID = repo.lastID
	repo.lastID++
	repo.data = append(repo.data, worker)
	return worker.ID, nil
}

// GetByID gets worker by ID.
// BUG: If user with this ID doesn't exist panic can be raised
func (repo *workerRepo) GetByID(id uint64) (*models.Worker, error) {
	if id > repo.lastID {
		// should return error!!
		// remade with own_errors
		return nil, nil
	}
	return repo.data[id], nil
}

func (repo *workerRepo) GetByEmail(email string) (*models.Worker, error) {
	for _, worker := range repo.data {
		if worker.Email == email {
			return worker, nil
		}
	}
	return nil, worker.ErrNoUserExist
}
