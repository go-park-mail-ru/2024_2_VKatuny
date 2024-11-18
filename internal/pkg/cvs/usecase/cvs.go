package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
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
		return nil, fmt.Errorf(dto.MsgDataBaseError)
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

func (cu *CVsUsecase) CreateCV(cv *dto.JSONCv, currentUser *dto.UserFromSession) (*dto.JSONCv, error) {
	cv.ApplicantID = currentUser.ID
	cv, err := cu.cvsRepo.Create(cv)
	if err != nil {
		cu.logger.Errorf("while adding to db got err: %s", err)
		return nil, err
	}
	return cv, nil
}

func (cu *CVsUsecase) GetCV(cvID uint64) (*dto.JSONCv, error) {
	cv, err := cu.cvsRepo.GetByID(cvID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return cv, nil
}

func (cu *CVsUsecase) UpdateCV(ID uint64, currentUser *dto.UserFromSession, cv *dto.JSONCv) (*dto.JSONCv, error) {
	oldCv, err := cu.cvsRepo.GetByID(ID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	if currentUser.ID != oldCv.ApplicantID {
		cu.logger.Errorf("not an owner tried to update CV, got %d expected %d", currentUser.ID, cv.ApplicantID)
		return nil, fmt.Errorf(dto.MsgAccessDenied)
	}
	cv.ApplicantID = oldCv.ApplicantID
	newCV, err := cu.cvsRepo.Update(ID, cv)

	if err != nil {
		cu.logger.Errorf("while updating in db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return newCV, nil
}

func (cu *CVsUsecase) DeleteCV(cvID uint64, currentUser *dto.UserFromSession) error {
	cv, err := cu.cvsRepo.GetByID(cvID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	if cv.ApplicantID != currentUser.ID {
		cu.logger.Errorf("not an owner tried to delete CV, got %d expected %d", currentUser.ID, cv.ApplicantID)
		return fmt.Errorf(dto.MsgAccessDenied)
	}
	err = cu.cvsRepo.Delete(cvID)
	if err != nil {
		cu.logger.Errorf("while deleting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	return nil
}
