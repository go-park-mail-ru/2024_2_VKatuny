package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	employerRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	"github.com/sirupsen/logrus"
)

type IEmployerUsecase interface {
	GetEmployerProfile(employerID uint64) (*dto.JSONGetEmployerProfile, error)
	UpdateEmployerProfile(employerID uint64, employerProfile *dto.JSONUpdateEmployerProfile) error
}

type EmployerUsecase struct {
	logger             *logrus.Logger
	employerRepository employerRepository.EmployerRepository
}

func NewEmployerUsecase(logger *logrus.Logger, repositories *internal.Repositories) *EmployerUsecase {
	employerRepository, ok := repositories.EmployerRepository.(employerRepository.EmployerRepository)
	if !ok {
		return nil
	}
	return &EmployerUsecase{
		logger:             logger,
		employerRepository: employerRepository,
	}
}

func (eu *EmployerUsecase) GetEmployerProfile(employerID uint64) (*dto.JSONGetEmployerProfile, error) {
	fn := "EmployerUsecase.GetEmployerProfile"
	employerModel, err := eu.employerRepository.GetByID(employerID)
	if err != nil {
		eu.logger.Debugf("function: %s; unable to get employer profile: %s", fn, err)
		return nil, err
	}
	eu.logger.Debugf("function: %s; got employer profile: %v", fn, employerModel)
	return &dto.JSONGetEmployerProfile{
		FirstName:          employerModel.FirstName,
		LastName:           employerModel.LastName,
		City:               employerModel.CityName,
		Position:           employerModel.Position,
		Company:            employerModel.CompanyName,
		CompanyDescription: employerModel.CompanyDescription,
		CompanyWebsite:     employerModel.CompanyWebsite,
		Avatar:             employerModel.PathToProfileAvatar,
		Contacts:           employerModel.Contacts,
	}, nil
}

func (eu *EmployerUsecase) UpdateEmployerProfile(employerID uint64, employerProfile *dto.JSONUpdateEmployerProfile) error {
	fn := "EmployerUsecase.UpdateEmployerProfile"
	err := eu.employerRepository.Update(employerID, employerProfile)
	if err != nil {
		eu.logger.Errorf("function: %s; unable to update employer profile: %s", fn, err)
		return err
	}
	return nil
}
