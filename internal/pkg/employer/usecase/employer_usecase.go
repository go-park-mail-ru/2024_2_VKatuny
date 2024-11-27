// Package usecase contains usecase for employer
package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// CreateEmployer accepts employer repository and validated form and creates new employer
func (u *EmployerUsecase) Create(ctx context.Context, form *dto.JSONEmployerRegistrationForm) (*dto.JSONUser, error) {
	fn := "EmployerUsecase.Create"

	u.logger = utils.SetLoggerRequestID(ctx, u.logger)
	u.logger.Debugf("%s: entering", fn)

	_, err := u.employerRepository.GetByEmail(form.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgUserAlreadyExists)
	}
	u.logger.Debugf("%s: OK, such user doesn't exist", fn)

	employer := &dto.EmployerInput{
		FirstName:          form.FirstName,
		LastName:           form.LastName,
		Position:           form.Position,
		CompanyName:        form.Company,
		CompanyDescription: form.CompanyDescription,
		CompanyWebsite:     form.CompanyWebsite,
		Email:              form.Email,
		Password:           utils.HashPassword(form.Password),
	}

	employerModel, err := u.employerRepository.Create(employer)
	if err != nil {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	u.logger.Debugf("%s: user created", fn)

	return &dto.JSONUser{
		ID:       employerModel.ID,
		UserType: dto.UserTypeApplicant,
	}, nil
}

func (u *EmployerUsecase) GetByID(ctx context.Context, ID uint64) (*dto.JSONEmployer, error) {
	fn := "EmployerUsecase.GetByID"

	u.logger = utils.SetLoggerRequestID(ctx, u.logger)
	u.logger.Debugf("%s: entering", fn)

	employerModel, err := u.employerRepository.GetByID(ID)
	if err != nil {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	u.logger.Debugf("%s: successfully got user", fn)

	return &dto.JSONEmployer{
		UserType:            dto.UserTypeEmployer,
		ID:                  employerModel.ID,
		FirstName:           employerModel.FirstName,
		LastName:            employerModel.LastName,
		CityName:            employerModel.CityName,
		Position:            employerModel.Position,
		CompanyName:         employerModel.CompanyName,
		CompanyDescription:  employerModel.CompanyDescription,
		CompanyWebsite:      employerModel.CompanyWebsite,
		Contacts:            employerModel.Contacts,
		PathToProfileAvatar: employerModel.PathToProfileAvatar,
		Email:               employerModel.Email,
		CreatedAt:           employerModel.CreatedAt,
		UpdatedAt:           employerModel.UpdatedAt,
		CompressedAvatar:    employerModel.CompressedAvatar,
	}, nil
}
