// Package usecase contains usecase for employer
package usecase

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	repoSession "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
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

func CreateEmployer(repo repository.EmployerRepository, sessionRepoApplicant repoSession.SessionRepository, form *dto.EmployerInput) (*dto.EmployerOutput, string, error) {
	form.Password = utils.HashPassword(form.Password)
	user, err := repo.Create(form)
	if err != nil {
		return nil, "", err
	}
	sessionID := utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeEmployer)
	sessionRepoApplicant.Create(user.ID, sessionID)
	return &dto.EmployerOutput{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		CityName:            user.CityName,
		Position:            user.Position,
		CompanyName:         user.CompanyName,
		CompanyDescription:  user.CompanyDescription,
		CompanyWebsite:      user.CompanyWebsite,
		PathToProfileAvatar: user.PathToProfileAvatar,
		Contacts:            user.Contacts,
		Email:               user.Email,
		PasswordHash:        user.PasswordHash,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, sessionID, err
}
