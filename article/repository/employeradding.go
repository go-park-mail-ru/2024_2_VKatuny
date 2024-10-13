package repository

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
)

func CreateEmployer(h *BD.EmployerHandlers, newUserInput *BD.EmployerInput) (BD.Employer, error) {
	_, err := GetEmployerByEmail(newUserInput.EmployerEmail)
	if err == nil {
		return BD.Employer{}, fmt.Errorf("User exist")
	} else {
		hashed := service.HashPassword(newUserInput.EmployerPassword)
		var id uint64 = h.Amount + 1
		h.Mu.Lock()
		h.Amount += 1
		h.Users[newUserInput.EmployerEmail] = BD.Employer{
			ID:                 id,
			EmployerName:       newUserInput.EmployerName,
			EmployerLastName:   newUserInput.EmployerLastName,
			EmployerPosition:   newUserInput.EmployerPosition,
			CompanyName:        newUserInput.CompanyName,
			CompanyDescription: newUserInput.CompanyDescription,
			Website:            newUserInput.Website,
			EmployerEmail:      newUserInput.EmployerEmail,
			EmployerPassword:   hashed,
		}
		h.Mu.Unlock()
		log.Println("employer registrated")
		log.Println(h.Users)
		return h.Users[newUserInput.EmployerEmail], nil
	}
}
