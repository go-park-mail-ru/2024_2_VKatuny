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

// GenerateSessionTokenWithLength generate random string with given length for session id
func GenerateSessionToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var ErrEmptyCookie = fmt.Errorf("client have an empty cookie")

func CheckAuthorization(session *http.Cookie, repoApplicant, repoEmployer session.Repository) (*dto.UserWithSession, error) {
	if sessionID := session.Value; sessionID != "" {

		id, err := repoApplicant.GetUserIdBySession(sessionID)
		userType := dto.UserTypeApplicant

		if err != nil {
			id, err = repoEmployer.GetUserIdBySession(sessionID)
			userType = dto.UserTypeEmployer
		}
		return &dto.UserWithSession{ID: id, UserType: userType, SessionID: sessionID}, err
	}
	return nil, ErrEmptyCookie
}

var (
	// ErrWrongPassword means that password is wrong
	ErrWrongPassword = fmt.Errorf("wrong password")
	// ErrNoApplicantWithSuchEmail means that there is no applicant with such email
	ErrNoApplicantWithSuchEmail = fmt.Errorf("there is no applicant with such email")
	// ErrNoEmployerWithSuchEmail means that there is no employer with such email
	ErrNoEmployerWithSuchEmail = fmt.Errorf("there is no employer with such email")
)

// LoginValidate ! TODO: rename function to more accurate meaning
func LoginValidate(newUserInput *dto.JSONLoginForm, repoApplicant worker.Repository, repoEmployer repository.EmployerRepository) (user *dto.UserIDAndType, err error) {
	if newUserInput.UserType == dto.UserTypeApplicant {
		var worker *models.Worker
		worker, err = repoApplicant.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, ErrNoApplicantWithSuchEmail
		}
		if !utils.EqualHashedPasswords(worker.Password, newUserInput.Password) {
			return nil, ErrWrongPassword
		}
		user = &dto.UserIDAndType{ID: worker.ID, UserType: dto.UserTypeApplicant}
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		var employer *models.Employer
		employer, err = repoEmployer.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, ErrNoEmployerWithSuchEmail
		}
		if !utils.EqualHashedPasswords(employer.Password, newUserInput.Password) {
			return nil, ErrWrongPassword
		}
		user = &dto.UserIDAndType{ID: employer.ID, UserType: dto.UserTypeEmployer}
	}
	return user, err
}

// LogoutValidate tries to remove session from db
// TODO: rename function to more accurate meaning
func LogoutValidate(newUserInput *dto.JSONLogoutForm, sessionID string, repoApplicant, repoEmployer session.Repository) error {
	var err error
	if newUserInput.UserType == dto.UserTypeApplicant {
		err = repoApplicant.Delete(sessionID)
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		err = repoEmployer.Delete(sessionID)
	}
	return err
}

func AddSession(repoApplicant, repoEmployer session.Repository, user *dto.UserIDAndType) (string, error) {
	sessionID := GenerateSessionToken(tokenLength)
	if user.UserType == dto.UserTypeApplicant {
		return sessionID, repoApplicant.Add(user.ID, sessionID)
	} else if user.UserType == dto.UserTypeEmployer {
		return sessionID, repoEmployer.Add(user.ID, sessionID)
	}
	return sessionID, nil
}
