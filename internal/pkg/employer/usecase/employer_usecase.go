// Package usecase contains usecase for employer
package usecase

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// CreateEmployerInputCheck accepts registartion form and checks it
func CreateEmployerInputCheck(form *dto.EmployerInput) error {
	if len(form.FirstName) > 50 || len(form.LastName) > 50 || len(form.Position) > 50 ||
		len(form.CompanyName) > 50 || strings.Index(form.Email, "@") < 0 || len(form.Password) > 50 {
		return fmt.Errorf("employer's fields aren't valid %s %s %s %s %s",
			form.FirstName,
			form.LastName,
			form.Position,
			form.CompanyName,
			form.Email,
		)
	}
	return nil
}

// CreateEmployer accepts employer repository and validated form and creates new employer
func (u *EmployerUsecase) Create(form *dto.JSONEmployerRegistrationForm) (*dto.JSONUser, error) {
	fn := "EmployerUsecase.Create"
	u.logger.Debugf("%s: entering", fn)

	_, err := u.employerRepository.GetByEmail(form.Email)
	if err.Error() != "sql: no rows in result set" {
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

func (u *EmployerUsecase) GetByID(ID uint64) (*dto.JSONEmployer, error) {
	fn := "EmployerUsecase.GetByID"

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
	}, nil
}
