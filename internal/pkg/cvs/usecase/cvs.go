package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type CVsUsecase struct {
	logger      *logrus.Entry
	cvsRepo     cvs.ICVsRepository
}

func NewCVsUsecase(logger *logrus.Logger, repositories *internal.Repositories) *CVsUsecase {
	return &CVsUsecase{
		logger:      logrus.NewEntry(logger),
		cvsRepo:     repositories.CVRepository,
	}
}

const (
	defaultVacanciesOffset = 0
	defaultVacanciesNum    = 10
)

var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

func (u *CVsUsecase) ValidateQueryParameters(ctx context.Context, offsetStr, numStr string) (uint64, uint64, error) {
	fn := "VacanciesUsecase.ValidateQueryParameters"
	u.logger = utils.SetLoggerRequestID(ctx, u.logger)

	var err error
	offset, err1 := strconv.Atoi(offsetStr)

	if err1 != nil {
		u.logger.Errorf("%s: query parameter offset isn't a number: %s", fn, err1)
		offset = 0
		err = ErrOffsetIsNotANumber
	}

	num, err2 := strconv.Atoi(numStr)
	if err2 != nil {
		u.logger.Errorf("%s:query parameter num isn't a number: %s", fn, err2)
		num = 0
		err = ErrNumIsNotANumber // previous err will be overwritten
	}
	return uint64(offset), uint64(num), err
}

func (cu *CVsUsecase) SearchCVs(ctx context.Context, offsetStr, numStr, searchStr, group, searchBy string) ([]*dto.JSONGetApplicantCV, error) {
	fn := "VacanciesUsecase.GetVacanciesWithOffset"
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	offset, num, err := cu.ValidateQueryParameters(ctx, offsetStr, numStr)
	if errors.Is(ErrOffsetIsNotANumber, err) {
		offset = defaultVacanciesOffset
	}
	if errors.Is(ErrNumIsNotANumber, err) {
		num = defaultVacanciesNum
	}
	if searchStr == "" && searchBy != "" {
		searchBy = ""
	}
	var CVsModel []*dto.JSONCv

	CVsModel, err = cu.cvsRepo.SearchAll(ctx, offset, offset+num, searchStr, group, searchBy)
	var CVs []*dto.JSONGetApplicantCV
	for _, CVModel := range CVsModel {
		CVs = append(CVs, &dto.JSONGetApplicantCV{
			ID:                   CVModel.ID,
			ApplicantID:          CVModel.ApplicantID,
			PositionRu:           CVModel.PositionRu,
			PositionEn:           CVModel.PositionEn,
			JobSearchStatus:      CVModel.JobSearchStatusName,
			WorkingExperience:    CVModel.WorkingExperience,
			PositionCategoryName: CVModel.PositionCategoryName,
			Avatar:               CVModel.Avatar,
			CreatedAt:            CVModel.CreatedAt,
			UpdatedAt:            CVModel.UpdatedAt,
			Description:          CVModel.Description,
			CompressedAvatar:     CVModel.CompressedAvatar,
		})
	}
	if err != nil {
		return nil, err
	}
	cu.logger.Debugf("%s: got %d vacancies", fn, len(CVs))
	return CVs, nil
}

func (cu *CVsUsecase) GetApplicantCVs(ctx context.Context, applicantID uint64) ([]*dto.JSONGetApplicantCV, error) {
	fn := "CVsUsecase.GetApplicantCVs"
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	CVsModel, err := cu.cvsRepo.GetCVsByApplicantID(ctx, applicantID)
	if err != nil {
		cu.logger.Errorf("function %s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	cu.logger.Debugf("function %s: success, got CVs from repository: %d", fn, len(CVsModel))

	CVs := make([]*dto.JSONGetApplicantCV, 0, len(CVsModel))
	for _, CVModel := range CVsModel {
		CVs = append(CVs, &dto.JSONGetApplicantCV{
			ID:                   CVModel.ID,
			ApplicantID:          CVModel.ApplicantID,
			PositionRu:           CVModel.PositionRu,
			PositionEn:           CVModel.PositionEn,
			JobSearchStatus:      CVModel.JobSearchStatusName,
			WorkingExperience:    CVModel.WorkingExperience,
			PositionCategoryName: CVModel.PositionCategoryName,
			Avatar:               CVModel.Avatar,
			CreatedAt:            CVModel.CreatedAt,
			UpdatedAt:            CVModel.UpdatedAt,
			Description:          CVModel.Description,
			CompressedAvatar:     CVModel.CompressedAvatar,
		})
	}

	return CVs, nil
}

func (cu *CVsUsecase) CreateCV(ctx context.Context, cv *dto.JSONCv, currentUser *dto.UserFromSession) (*dto.JSONCv, error) {
	fn := "CVsUsecase.CreateCV"
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	cv.ApplicantID = currentUser.ID
	cv, err := cu.cvsRepo.Create(ctx, cv)
	if err != nil {
		cu.logger.Errorf("%s: while adding to db got err: %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return cv, nil
}

func (cu *CVsUsecase) GetCV(ctx context.Context, cvID uint64) (*dto.JSONCv, error) {
	fn := "CVsUsecase.CreateCV"
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	cv, err := cu.cvsRepo.GetByID(ctx, cvID)
	if err != nil {
		cu.logger.Errorf("%s: while getting from db got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return cv, nil
}

func (cu *CVsUsecase) UpdateCV(ctx context.Context, ID uint64, currentUser *dto.UserFromSession, cv *dto.JSONCv) (*dto.JSONCv, error) {
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	oldCv, err := cu.cvsRepo.GetByID(ctx, ID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	if currentUser.ID != oldCv.ApplicantID {
		cu.logger.Errorf("not an owner tried to update CV, got %d expected %d", currentUser.ID, cv.ApplicantID)
		return nil, fmt.Errorf(dto.MsgAccessDenied)
	}
	cv.ApplicantID = oldCv.ApplicantID
	newCV, err := cu.cvsRepo.Update(ctx, ID, cv)

	if err != nil {
		cu.logger.Errorf("while updating in db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return newCV, nil
}

func (cu *CVsUsecase) DeleteCV(ctx context.Context, cvID uint64, currentUser *dto.UserFromSession) error {
	cu.logger = utils.SetLoggerRequestID(ctx, cu.logger)

	cv, err := cu.cvsRepo.GetByID(ctx, cvID)
	if err != nil {
		cu.logger.Errorf("while getting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	if cv.ApplicantID != currentUser.ID {
		cu.logger.Errorf("not an owner tried to delete CV, got %d expected %d", currentUser.ID, cv.ApplicantID)
		return fmt.Errorf(dto.MsgAccessDenied)
	}
	err = cu.cvsRepo.Delete(ctx, cvID)
	if err != nil {
		cu.logger.Errorf("while deleting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	return nil
}
