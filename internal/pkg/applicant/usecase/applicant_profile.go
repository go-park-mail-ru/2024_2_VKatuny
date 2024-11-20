package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	microservicesinterface "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/microservices_interfaces"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/sirupsen/logrus"
)

type ApplicantUsecase struct {
	logger          *logrus.Logger
	applicantRepo   applicant.IApplicantRepository
	compressManager microservicesinterface.ICompress
}

func NewApplicantUsecase(logger *logrus.Logger, repositories *internal.Repositories, CompressManager microservicesinterface.ICompress) *ApplicantUsecase {
	return &ApplicantUsecase{
		logger:          logger,
		applicantRepo:   repositories.ApplicantRepository,
		compressManager: CompressManager,
	}
}

// GetApplicantProfile accepts the profile of an applicant using the given userID.
// It logs the process of fetching the profile and returns the applicant profile data
// encapsulated in a JSONGetApplicantProfile DTO. If fetching the profile fails,
// it returns an error.
func (au *ApplicantUsecase) GetApplicantProfile(userID uint64) (*dto.JSONGetApplicantProfile, error) {
	fn := "ApplicantUsecase.GetApplicantProfile"

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
	}, nil
}

// UpdateApplicantProfile updates the profile of an applicant with the given ID
// using the provided new profile data. It logs the update process and returns
// an error if the update fails.
func (au *ApplicantUsecase) UpdateApplicantProfile(applicantID uint64, newProfileData *dto.JSONUpdateApplicantProfile) error {
	fn := "ApplicantUsecase.UpdateApplicantProfile"

	au.logger.Debugf("function: %s; applicant id: %d. Trying to update applicant profile", fn, applicantID)

	_, err := au.applicantRepo.Update(applicantID, newProfileData)
	if err != nil {
		au.logger.Errorf("function: %s; got err: %s", fn, err)
		return err
	}
	_, err = au.compressManager.CompressAndSaveFile(
		context.Background(),
		&compressmicroservice.CompressAndSaveFileInput{
			FileName: "filename",
			FileType: "filetype",
		})
	if err != nil {
		au.logger.Errorf("fail compress microservice")
		return err
	}

	au.logger.Debugf("function: %s; successfully updated applicant profile", fn)
	return nil
}
