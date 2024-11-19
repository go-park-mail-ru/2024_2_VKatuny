package cvs

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// Здесь все интерфейсы взаимодействия с CVs

type ICVsRepository interface {
	GetCVsByApplicantID(applicantID uint64) ([]*dto.JSONCv, error)
	Create(cv *dto.JSONCv) (*dto.JSONCv, error)
	GetByID(ID uint64) (*dto.JSONCv, error)
	Update(ID uint64, updatedCv *dto.JSONCv) (*dto.JSONCv, error)
	Delete(ID uint64) error
	GetWithOffset(offset uint64, num uint64) ([]*dto.JSONCv, error)
	SearchByPositionDescription(offset uint64, num uint64, searchStr string) ([]*dto.JSONCv, error)
}

type ICVsUsecase interface {
	GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error)
	CreateCV(cv *dto.JSONCv, currentUser *dto.UserFromSession) (*dto.JSONCv, error)
	GetCV(ID uint64) (*dto.JSONCv, error)
	UpdateCV(ID uint64, currentUser *dto.UserFromSession, cv *dto.JSONCv) (*dto.JSONCv, error)
	DeleteCV(ID uint64, currentUser *dto.UserFromSession) error
	SearchCVs(offsetStr, numStr, searchStr string) ([]*dto.JSONGetApplicantCV, error)
}
