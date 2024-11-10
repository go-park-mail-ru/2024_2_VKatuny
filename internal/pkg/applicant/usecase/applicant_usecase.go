// Package usecase contains usecase for worker
package usecase

import (
	"fmt"
	"strings"

	repoApplicant "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	repoSession "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
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

// CreateApplicant accepts employer repository and validated form and creates new employer
func CreateApplicant(repo repoApplicant.IApplicantRepository, sessionRepoApplicant repoSession.SessionRepository, form *dto.ApplicantInput) (*dto.ApplicantOutput, string, error) {
	_, err := repo.GetByEmail(form.Email)
	// if applicant != nil {
	// 	return nil, "", fmt.Errorf(dto.MsgUserAlreadyExists)
	// }
	// if err != nil {
	// 	return nil, "", fmt.Errorf(dto.MsgDataBaseError)
	// }
	if err.Error() != "sql: no rows in result set" {
		return nil, "", fmt.Errorf(dto.MsgDataBaseError)
	}
	form.Password = utils.HashPassword(form.Password)
	user, err := repo.Create(form)
	if err != nil {
		return nil, "", fmt.Errorf(dto.MsgDataBaseError)
	}
	sessionID := utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeApplicant)
	err = sessionRepoApplicant.Create(user.ID, sessionID)
	// if err != nil {
	// 	return nil, "", fmt.Errorf(dto.MsgDataBaseError)
	// }
	return &dto.ApplicantOutput{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		CityName:            user.CityName,
		BirthDate:           user.BirthDate,
		PathToProfileAvatar: user.PathToProfileAvatar,
		Contacts:            user.Contacts,
		Education:           user.Education,
		Email:               user.Email,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, sessionID, err // TODO: return nil instead of err
}
