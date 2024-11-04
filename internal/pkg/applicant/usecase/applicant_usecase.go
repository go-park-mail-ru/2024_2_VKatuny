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

func CreateApplicantInputCheck(Name, LastName, Email, Password string) error {
	if len(Name) > 50 || len(LastName) > 50 ||
		strings.Index(Email, "@") < 0 || len(Password) > 50 {
		return fmt.Errorf("applicant's fields aren't valid %s %s %s", Name, LastName, Email)
	}
	return nil
}

// CreateApplicant accepts employer repository and validated form and creates new employer
func CreateApplicant(repo repoApplicant.ApplicantRepository, sessionRepoApplicant repoSession.SessionRepository, form *dto.ApplicantInput) (*dto.ApplicantOutput, string, error) {
	form.Password = utils.HashPassword(form.Password)
	user, err := repo.Create(form)
	if err != nil {
		return nil, "", err
	}
	sessionID := utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeApplicant)
	sessionRepoApplicant.Create(user.ID, sessionID)
	return &dto.ApplicantOutput{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		CityName:            user.CityName,
		BirthDate:           user.BirthDate,
		PathToProfileAvatar: user.LastName,
		Constants:           user.Contacts,
		Education:           user.Education,
		Email:               user.Email,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, sessionID, err
}
