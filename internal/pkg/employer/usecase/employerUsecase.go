// Package usecase contains usecase for employer
package usecase

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

// CreateEmployerInputCheck should be remade
func CreateEmployerInputCheck(form *dto.JSONEmployerRegistrationForm) error {
	if len(form.Name) > 50 || len(form.LastName) > 50 || len(form.Position) > 50 ||
		len(form.CompanyName) > 50 || strings.Index(form.Email, "@") < 0 || len(form.Password) > 50 {
		return fmt.Errorf("employer's fields aren't valid %s %s %s %s %s",
		       form.Name,
			   form.LastName,
			   form.Position,
			   form.CompanyName,
			   form.Email,
			)
	}
	return nil
}
