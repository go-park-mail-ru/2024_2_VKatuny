package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/sirupsen/logrus"
)

type CVsUsecase struct {
	logger      *logrus.Logger
	cvsRepo     cvs.ICVsRepository
	sessionRepo session.ISessionRepository
}

func NewCVsUsecase(logger *logrus.Logger, repositories *internal.Repositories) *CVsUsecase {
	return &CVsUsecase{
		logger:      logger,
		cvsRepo:     repositories.CVRepository,
		sessionRepo: repositories.SessionApplicantRepository,
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
			PositionRu:        CVModel.PositionRus,
			PositionEn:        CVModel.PositionEng,
			JobSearchStatus:   CVModel.JobSearchStatus,
			WorkingExperience: CVModel.WorkingExperience,
			CreatedAt:         CVModel.CreatedAt.Format("2006.01.02 15:02:39"),
		})
	}

	return CVs, nil
}

func (cu *CVsUsecase) CreateCV(cv *dto.JSONCv) (*dto.JSONCv, error) {
	newCVModel, err := cu.cvsRepo.Create(cv)
	if err != nil {
		cu.logger.Errorf("while adding to db got err: %s", err)
		return nil, err
	}
	return &dto.JSONCv{
		ApplicantID:       newCVModel.ApplicantID,
		PositionRu:        newCVModel.PositionRus,
		PositionEn:        newCVModel.PositionEng,
		Description:       newCVModel.Description,
		JobSearchStatus:   newCVModel.JobSearchStatus,
		WorkingExperience: newCVModel.WorkingExperience,
	}, nil
}

func (cu *CVsUsecase) GetCV(cvID uint64) (*dto.JSONCv, error) {
	CVModel, err := cu.cvsRepo.GetByID(cvID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return nil, err
	}
	return &dto.JSONCv{
		ApplicantID:       CVModel.ApplicantID,
		PositionRu:        CVModel.PositionRus,
		PositionEn:        CVModel.PositionEng,
		Description:       CVModel.Description,
		JobSearchStatus:   CVModel.JobSearchStatus,
		WorkingExperience: CVModel.WorkingExperience,
	}, nil
}

func (cu *CVsUsecase) UpdateCV(ID uint64 , sessionID string, cv *dto.JSONCv) (*dto.JSONCv, error) {
	currentUserID, err := cu.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return nil, commonerrors.ErrSessionNotFound
	}

	if currentUserID != cv.ApplicantID {
		cu.logger.Errorf("not an owner tried to update CV, got %d expected %d", currentUserID, cv.ApplicantID)
		return nil, commonerrors.ErrUnauthorized
	}

	CVModel, err := cu.cvsRepo.Update(ID, cv)
	if err != nil {
		cu.logger.Errorf("while updating in db got err %s", err)
		return nil, err
	}
	return &dto.JSONCv{
		ApplicantID:       CVModel.ApplicantID,
		PositionRu:        CVModel.PositionRus,
		PositionEn:        CVModel.PositionEng,
		Description:       CVModel.Description,
		JobSearchStatus:   CVModel.JobSearchStatus,
		WorkingExperience: CVModel.WorkingExperience,
	}, nil
}

func (cu *CVsUsecase) DeleteCV(cvID uint64, sessionID string) error {
	currentUserID, err := cu.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		cu.logger.Errorf("while getting sesion from db got err %s", err)
		return commonerrors.ErrSessionNotFound
	}

	CVModel, err := cu.cvsRepo.GetByID(cvID)
	if CVModel.ApplicantID != currentUserID {
		cu.logger.Errorf("not an owner tried to delete CV, got %d expected %d", currentUserID, CVModel.ApplicantID)
		return commonerrors.ErrUnauthorized
	}
	err = cu.cvsRepo.Delete(cvID)
	if err != nil {
		cu.logger.Errorf("while deleting from db got err %s", err)
		return err
	}
	return nil
}
