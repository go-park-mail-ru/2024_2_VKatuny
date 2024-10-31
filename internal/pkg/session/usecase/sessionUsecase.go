package sessionUsecase

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker"
	"golang.org/x/crypto/bcrypt"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const tokenLength = 32

// RandStringRunes generate random string with length of n for session id
func GenerateSessionTokenWithLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// HashPassword hashing password with cost og 10
func HashPassword(password string) string {
	bytePassword := []byte(password)
	cost := 10
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	return string(hashedPassword[:])
}

// EqualHashedPasswords returns true if first argument was hashed from second
func EqualHashedPasswords(passwordBD string, passwordFront string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordBD), []byte(passwordFront))
	return err == nil
}

var ErrEmptyCookie = fmt.Errorf("client have an empty cookie")

func SessionCheck(session *http.Cookie, repoApplicant, repoEmployer session.Repository) (uint64, string, string, error) {
	if sessionId := session.Value; sessionId != "" {

		id, err := repoApplicant.GetUserIdBySession(sessionId)
		userType := dto.UserTypeApplicant

		if err != nil {
			id, err = repoEmployer.GetUserIdBySession(sessionId)
			userType = dto.UserTypeEmployer
		}
		return id, userType, sessionId, err
	} else {
		return 0, "", "", ErrEmptyCookie
	}

}

var ErrWrongPassword = fmt.Errorf("wrong password")
var ErrNoApplicantWithSuchEmail = fmt.Errorf("there is no applicant with such email")
var ErrNoEmployerWithSuchEmail = fmt.Errorf("there is no employer with such email")

func LoginCheck(newUserInput *dto.JsonLoginForm, repoApplicant worker.Repository, repoApplicantSession, repoEmployerSession session.Repository, repoEmployer employer.Repository) (string, error) {
	sessionId := GenerateSessionTokenWithLength(tokenLength)
	var err error
	if newUserInput.UserType == dto.UserTypeApplicant {
		var user *models.Worker
		user, err = repoApplicant.GetByEmail(newUserInput.Email)
		if err != nil {
			return "", ErrNoApplicantWithSuchEmail
		}
		if !EqualHashedPasswords(user.Password, newUserInput.Password) {
			return "", ErrWrongPassword
		}
		err = repoApplicantSession.Add(user.ID, sessionId)
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		var user *models.Employer
		user, err = repoEmployer.GetByEmail(newUserInput.Email)
		if err != nil {
			return "", ErrNoEmployerWithSuchEmail
		}
		if !EqualHashedPasswords(user.Password, newUserInput.Password) {
			return "", ErrWrongPassword
		}
		err = repoEmployerSession.Add(user.ID, sessionId)
	}
	return sessionId, err
}

func LogoutCheck(newUserInput *dto.JsonLogoutForm, sessionId string, repoApplicant, repoEmployer session.Repository) error {
	var err error
	if newUserInput.UserType == dto.UserTypeApplicant {
		err = repoApplicant.Delete(sessionId)
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		err = repoEmployer.Delete(sessionId)
	}
	return err
}
