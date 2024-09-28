package service

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func TryCreateEmployer(h *BD.EmployerHandlers, newUserInput *BD.EmployerInput) (BD.Employer, error) {
	_, rErr := storage.GetEmployerByEmail(h, newUserInput.EmployerEmail)
	fmt.Println(rErr)
	if rErr != nil {
		h.Mu.Lock()
		var id uint64 = h.Amount + 1
		h.Users[newUserInput.EmployerEmail] = BD.Employer{
			ID:                 id,
			EmployerName:       newUserInput.EmployerName,
			EmployerLastName:   newUserInput.EmployerLastName,
			EmployerPosition:   newUserInput.EmployerPosition,
			CompanyName:        newUserInput.CompanyName,
			CompanyDescription: newUserInput.CompanyDescription,
			Website:            newUserInput.Website,
			EmployerEmail:      newUserInput.EmployerEmail,
			EmployerPassword:   storage.HashPassword(newUserInput.EmployerPassword),
		}
		h.Mu.Unlock()
		return h.Users[newUserInput.EmployerEmail], nil
	} else {
		return BD.Employer{}, rErr
	}
}
