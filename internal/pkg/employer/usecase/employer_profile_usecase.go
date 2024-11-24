package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

type EmployerUsecase struct {
	logger             *logrus.Entry
	employerRepository employer.IEmployerRepository
}

func NewEmployerUsecase(logger *logrus.Logger, repositories *internal.Repositories) *EmployerUsecase {
	return &EmployerUsecase{
		logger:             &logrus.Entry{Logger: logger},
		employerRepository: repositories.EmployerRepository,
	}
}

func (eu *EmployerUsecase) GetEmployerProfile(ctx context.Context, employerID uint64) (*dto.JSONGetEmployerProfile, error) {
	fn := "EmployerUsecase.GetEmployerProfile"

	eu.logger = utils.SetLoggerRequestID(ctx, eu.logger)
	eu.logger.Debugf("%s: entering", fn)

	employerModel, err := eu.employerRepository.GetByID(employerID)
	if err != nil {
		eu.logger.Debugf("%s: unable to get employer profile: %s", fn, err)
		return nil, err
	}
	eu.logger.Debugf("%s: got employer profile: %v", fn, employerModel)
	return &dto.JSONGetEmployerProfile{
		ID:                 employerModel.ID,
		FirstName:          employerModel.FirstName,
		LastName:           employerModel.LastName,
		City:               employerModel.CityName,
		Position:           employerModel.Position,
		Company:            employerModel.CompanyName,
		CompanyDescription: employerModel.CompanyDescription,
		CompanyWebsite:     employerModel.CompanyWebsite,
		Contacts:           employerModel.Contacts,
		Avatar:             employerModel.PathToProfileAvatar,
	}, nil
}

func (eu *EmployerUsecase) UpdateEmployerProfile(ctx context.Context, employerID uint64, employerProfile *dto.JSONUpdateEmployerProfile) error {
	fn := "EmployerUsecase.UpdateEmployerProfile"

	eu.logger = utils.SetLoggerRequestID(ctx, eu.logger)
	eu.logger.Debugf("%s: entering", fn)

	_, err := eu.employerRepository.Update(employerID, employerProfile)
	if err != nil {
		eu.logger.Errorf("function: %s; unable to update employer profile: %s", fn, err)
		return err
	}
	return nil
}
