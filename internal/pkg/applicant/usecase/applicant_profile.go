package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	"github.com/sirupsen/logrus"
)

type ApplicantUsecase struct {
	logger        *logrus.Entry
	applicantRepo applicant.IApplicantRepository
}

func NewApplicantUsecase(logger *logrus.Logger, repositories *internal.Repositories) *ApplicantUsecase {
	return &ApplicantUsecase{
		logger:        logrus.NewEntry(logger),
		applicantRepo: repositories.ApplicantRepository,
	}
}

// GetApplicantProfile accepts the profile of an applicant using the given userID.
// It logs the process of fetching the profile and returns the applicant profile data
// encapsulated in a JSONGetApplicantProfile DTO. If fetching the profile fails,
// it returns an error.
func (au *ApplicantUsecase) GetApplicantProfile(ctx context.Context, userID uint64) (*dto.JSONGetApplicantProfile, error) {
	fn := "ApplicantUsecase.GetApplicantProfile"
	au.logger = utils.SetLoggerRequestID(ctx, au.logger)
	au.logger.Debugf("%s: entering", fn)

	au.logger.Debugf("function: %s; user id: %d. Trying to get applicant profile by id", fn, userID)
	applicantModel, err := au.applicantRepo.GetByID(userID)
	if err != nil {
		au.logger.Errorf("function: %s; got err: %s", fn, err)
		return nil, err
	}
	au.logger.Debugf("function: %s; successfully got applicant profile: %v", fn, applicantModel)
	return &dto.JSONGetApplicantProfile{
		ID:        applicantModel.ID,
		FirstName: applicantModel.FirstName,
		LastName:  applicantModel.LastName,
		City:      applicantModel.CityName,
		BirthDate: applicantModel.BirthDate,
		Contacts:  applicantModel.Contacts,
		Education: applicantModel.Education,
		Avatar:    applicantModel.PathToProfileAvatar,
		CompressedAvatar: applicantModel.CompressedAvatar,
	}, nil
}

// UpdateApplicantProfile updates the profile of an applicant with the given ID
// using the provided new profile data. It logs the update process and returns
// an error if the update fails.
func (au *ApplicantUsecase) UpdateApplicantProfile(ctx context.Context, applicantID uint64, newProfileData *dto.JSONUpdateApplicantProfile) error {
	fn := "ApplicantUsecase.UpdateApplicantProfile"
	au.logger = utils.SetLoggerRequestID(ctx, au.logger)
	au.logger.Debugf("%s: entering", fn)

	au.logger.Debugf("function: %s; applicant id: %d. Trying to update applicant profile", fn, applicantID)

	_, err := au.applicantRepo.Update(applicantID, newProfileData)
	if err != nil {
		au.logger.Errorf("function: %s; got err: %s", fn, err)
		return err
	}
	fmt.Println("compress")
	if err != nil {
		au.logger.Errorf("fail compress microservice")
		return err
	}

	au.logger.Debugf("function: %s; successfully updated applicant profile", fn)
	return nil
}
