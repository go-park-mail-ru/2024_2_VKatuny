package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

type ICVsUsecase interface {
	GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error)
}

type CVsUsecase struct {
	logger  *logrus.Logger
	cvsRepo cvs.ICVsRepository
}

func NewCVsUsecase(logger *logrus.Logger, repositories *internal.Repositories) *CVsUsecase {
	return &CVsUsecase{
		logger:  logger,
		cvsRepo: repositories.CVRepository,
	}
}

func (cu *CVsUsecase) GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error) {
	fn := "CVsUsecase.GetApplicantCVs"

	CVsModel, err := cu.cvsRepo.GetCVsByApplicantID(applicantID)
	if err != nil {
		cu.logger.Errorf("function %s: got err %s", fn, err)
		return nil, err
	}
	cu.logger.Debugf("function %s: success, got CVs from repository: %d", fn, len(CVsModel))

	CVs := make([]*dto.JSONGetApplicantCV, 0, len(CVsModel))
	for _, CVModel := range CVsModel {
		CVs = append(CVs, &dto.JSONGetApplicantCV{
			ID:                CVModel.ID,
			ApplicantID:       CVModel.ApplicantID,
			PositionRu:        CVModel.PositionRu,
			PositionEn:        CVModel.PositionEn,
			JobSearchStatus:   CVModel.JobSearchStatusName,
			WorkingExperience: CVModel.WorkingExperience,
			CreatedAt:         CVModel.CreatedAt,
		})
	}

	return CVs, nil
}
