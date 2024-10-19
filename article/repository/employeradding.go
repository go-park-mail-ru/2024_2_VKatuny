// Package repository is db interactions layer
package repository

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_2_VKatuny/article/usecase/service"
	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
)

// CreateEmployer creates employer in db
func CreateEmployer(h *inmemorydb.EmployerHandlers, newUserInput *inmemorydb.EmployerInput) (inmemorydb.Employer, error) {
	_, err := GetEmployerByEmail(newUserInput.EmployerEmail)
	if err == nil {
		return inmemorydb.Employer{}, fmt.Errorf("User exist")
	}
	hashed := service.HashPassword(newUserInput.EmployerPassword)
	var id uint64 = h.Amount + 1
	h.Mu.Lock()
	h.Amount++
	h.Users[newUserInput.EmployerEmail] = inmemorydb.Employer{
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
