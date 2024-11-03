// Package usecase contains usecase for worker
package usecase

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
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
func CreateApplicant(repo repository.ApplicantRepository, form *dto.ApplicantInput) (*models.Applicant, error) {
	form.Password = utils.HashPassword(form.Password)
	user, err := repo.Create(form)
	return user, err
}
