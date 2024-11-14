package usecase

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/sirupsen/logrus"
)

type EmployerUsecase struct {
	logger             *logrus.Entry
	employerRepository employer.IEmployerRepository
}

func NewEmployerUsecase(app *internal.App) *EmployerUsecase {
	return &EmployerUsecase{
		logger:             &logrus.Entry{Logger: app.Logger},
		employerRepository: app.Repositories.EmployerRepository,
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
		ID:                 employerModel.ID,
		FirstName:          employerModel.FirstName,
		LastName:           employerModel.LastName,
		City:               employerModel.CityName,
		Position:           employerModel.Position,
		Company:            employerModel.CompanyName,
		CompanyDescription: employerModel.CompanyDescription,
		CompanyWebsite:     employerModel.CompanyWebsite,
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
