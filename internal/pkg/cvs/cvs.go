package cvs

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// Здесь все интерфейсы взаимодействия с CVs

type ICVsRepository interface {
	GetCVsByApplicantID(ctx context.Context, applicantID uint64) ([]*dto.JSONCv, error)
	Create(ctx context.Context, cv *dto.JSONCv) (*dto.JSONCv, error)
	GetByID(ctx context.Context, ID uint64) (*dto.JSONCv, error)
	Update(ctx context.Context, ID uint64, updatedCv *dto.JSONCv) (*dto.JSONCv, error)
	Delete(ctx context.Context, ID uint64) error
	SearchAll(ctx context.Context, offset uint64, num uint64, searchStr, group, searchBy string) ([]*dto.JSONCv, error)
}

type ICVsUsecase interface {
	GetApplicantCVs(ctx context.Context, applicantID uint64) ([]*dto.JSONGetApplicantCV, error)
	CreateCV(ctx context.Context, cv *dto.JSONCv, currentUser *dto.UserFromSession) (*dto.JSONCv, error)
	GetCV(ctx context.Context, ID uint64) (*dto.JSONCv, error)
	UpdateCV(ctx context.Context, ID uint64, currentUser *dto.UserFromSession, cv *dto.JSONCv) (*dto.JSONCv, error)
	DeleteCV(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error
	SearchCVs(ctx context.Context, offsetStr, numStr, searchStr, group, searchBy string) ([]*dto.JSONGetApplicantCV, error)
}
