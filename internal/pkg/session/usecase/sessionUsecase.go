// Package usecase contains usecase for session
package usecase

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const tokenLength = 32

// GenerateSessionToken generate random string for session id with default length
func GenerateSessionToken() string {
	return GenerateSessionTokenWithLength(tokenLength)
}

// GenerateSessionTokenWithLength generate random string with given length for session id
func GenerateSessionTokenWithLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var ErrEmptyCookie = fmt.Errorf("client have an empty cookie")

func SessionCheck(session *http.Cookie, repoApplicant, repoEmployer session.Repository) (uint64, string, string, error) {
	if sessionID := session.Value; sessionID != "" {

		id, err := repoApplicant.GetUserIdBySession(sessionID)
		userType := dto.UserTypeApplicant

		if err != nil {
			id, err = repoEmployer.GetUserIdBySession(sessionID)
			userType = dto.UserTypeEmployer
		}
		return id, userType, sessionID, err
	}
	return 0, "", "", ErrEmptyCookie
}


var (
	// ErrWrongPassword means that password is wrong
	ErrWrongPassword = fmt.Errorf("wrong password")
	// ErrNoApplicantWithSuchEmail means that there is no applicant with such email
	ErrNoApplicantWithSuchEmail = fmt.Errorf("there is no applicant with such email")
	// ErrNoEmployerWithSuchEmail means that there is no employer with such email
	ErrNoEmployerWithSuchEmail = fmt.Errorf("there is no employer with such email")
)

// LoginCheck ! TODO: rename function to more accurate meaning
func LoginCheck(newUserInput *dto.JSONLoginForm, repoApplicant worker.Repository, repoEmployer repository.EmployerRepository) (uint64, error) {
	var err error
	var id uint64
	if newUserInput.UserType == dto.UserTypeApplicant {
		var user *models.Worker
		user, err = repoApplicant.GetByEmail(newUserInput.Email)
		if err != nil {
			return 0, ErrNoApplicantWithSuchEmail
		}
		if !utils.EqualHashedPasswords(user.Password, newUserInput.Password) {
			return 0, ErrWrongPassword
		}
		id = user.ID
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		var user *models.Employer
		user, err = repoEmployer.GetByEmail(newUserInput.Email)
		if err != nil {
			return 0, ErrNoEmployerWithSuchEmail
		}
		if !utils.EqualHashedPasswords(user.Password, newUserInput.Password) {
			return 0, ErrWrongPassword
		}
		id = user.ID
	}
	return id, err
}

// LogoutCheck tries to remove session from db
// TODO: rename function to more accurate meaning
func LogoutCheck(newUserInput *dto.JSONLogoutForm, sessionID string, repoApplicant, repoEmployer session.Repository) error {
	var err error
	if newUserInput.UserType == dto.UserTypeApplicant {
		err = repoApplicant.Delete(sessionID)
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		err = repoEmployer.Delete(sessionID)
	}
	return err
}
