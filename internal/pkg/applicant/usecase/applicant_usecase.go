// Package usecase contains usecase for worker
package usecase

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

// TODO: refactor valid code
func CreateApplicantInputCheck(Name, LastName, Email, Password string) error {
	if len(Name) > 50 || len(LastName) > 50 ||
		strings.Index(Email, "@") < 0 || len(Password) > 50 {
		return fmt.Errorf("applicant's fields aren't valid %s %s %s", Name, LastName, Email)
	}
	return nil
}

func (u *ApplicantUsecase) Create(form *dto.JSONApplicantRegistrationForm) (*dto.JSONUser, error) {
	fn := "ApplicantUsecase.Create"

	_, err := u.applicantRepo.GetByEmail(form.Email)
	if err.Error() != "sql: no rows in result set" {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgUserAlreadyExists)
	}
	u.logger.Debugf("%s: OK, such user doesn't exist", fn)

	applicant := &dto.ApplicantInput{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		BirthDate: form.BirthDate,
		Email:     form.Email,
		Password:  utils.HashPassword(form.Password),
	}
	
	applicantModel, err := u.applicantRepo.Create(applicant)
	if err != nil {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	u.logger.Debugf("%s: user created", fn)

	return &dto.JSONUser{
		ID:       applicantModel.ID,
		UserType: dto.UserTypeApplicant,
	}, nil
}

func (u *ApplicantUsecase) GetByID(ID uint64) (*dto.JSONApplicantOutput, error) {
	fn := "ApplicantUsecase.GetByID"
	u.logger.Debugf("%s: entering", fn)

	applicantModel, err := u.applicantRepo.GetByID(ID)
	if err != nil {
		u.logger.Errorf("%s: got err %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	u.logger.Debugf("%s: successfully got user", fn)

	return &dto.JSONApplicantOutput{
		ID:           applicantModel.ID,
		UserType:     dto.UserTypeApplicant,
		FirstName:    applicantModel.FirstName,
		LastName:     applicantModel.LastName,
		CityName:     applicantModel.CityName,
		BirthDate:    applicantModel.BirthDate,
		PathToProfileAvatar: applicantModel.PathToProfileAvatar,
		Contacts:     applicantModel.Contacts,
		Education:    applicantModel.Education,
		UpdatedAt:    applicantModel.UpdatedAt,
		CreatedAt:    applicantModel.CreatedAt,
	}, nil	
}
