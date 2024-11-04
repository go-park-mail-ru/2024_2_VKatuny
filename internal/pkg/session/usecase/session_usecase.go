// Package usecase contains usecase for session
package usecase

import (
	"fmt"
	"math/rand"
	"net/http"

	applicantRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	employerRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	sessionRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
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

func CheckAuthorization(newUserInput *dto.JSONLogoutForm, session *http.Cookie, sessionRepoApplicant sessionRepo.SessionRepository, sessionRepoEmployer sessionRepo.SessionRepository) (uint64, error) {
	if session == nil || session.Value == "" {
		return 0, ErrEmptyCookie
	}
	sessionID := session.Value
	if newUserInput.UserType == dto.UserTypeApplicant {
		id, err := sessionRepoApplicant.GetUserIdBySession(sessionID)
		return id, err
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		id, err := sessionRepoEmployer.GetUserIdBySession(sessionID)
		return id, err
	}
	return 0, fmt.Errorf("err type user")
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
func LoginValidate(newUserInput *dto.JSONLoginForm, repoApplicant applicantRepo.ApplicantRepository, repoEmployer employerRepo.EmployerRepository) (user *dto.UserIDAndType, err error) {
	if newUserInput.UserType == dto.UserTypeApplicant {
		var worker *models.Applicant
		worker, err = repoApplicant.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, ErrNoApplicantWithSuchEmail
		}
		if !utils.EqualHashedPasswords(worker.PasswordHash, newUserInput.Password) {
			return nil, ErrWrongPassword
		}
		user = &dto.UserIDAndType{ID: worker.ID, UserType: dto.UserTypeApplicant}
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		var employer *models.Employer
		employer, err = repoEmployer.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, ErrNoEmployerWithSuchEmail
		}
		if !utils.EqualHashedPasswords(employer.PasswordHash, newUserInput.Password) {
			return nil, ErrWrongPassword
		}
		user = &dto.UserIDAndType{ID: employer.ID, UserType: dto.UserTypeEmployer}
	}
	return user, err
}

// LogoutValidate tries to remove session from db
// TODO: rename function to more accurate meaning
func LogoutValidate(newUserInput *dto.JSONLogoutForm, session string, sessionRepoApplicant, sessionRepoEmployer sessionRepo.SessionRepository) (uint64, error) {
	if newUserInput.UserType == dto.UserTypeApplicant {
		id, err := sessionRepoApplicant.GetUserIdBySession(session)
		if err != nil {
			return 0, err
		}
		err = sessionRepoApplicant.Delete(session)
		if err == nil {
			return id, nil
		}
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		id, err := sessionRepoEmployer.GetUserIdBySession(session)
		if err != nil {
			return 0, err
		}
		err = sessionRepoEmployer.Delete(session)
		if err == nil {
			return id, nil
		}
	}
	return 0, fmt.Errorf("err type user")
}

func AddSession(sessionRepoApplicant, sessionRepoEmployer sessionRepo.SessionRepository, user *dto.UserIDAndType) (string, error) {
	sessionID := GenerateSessionToken(tokenLength)
	if user.UserType == dto.UserTypeApplicant {
		return sessionID, sessionRepoApplicant.Create(user.ID, sessionID)
	} else if user.UserType == dto.UserTypeEmployer {
		return sessionID, sessionRepoEmployer.Create(user.ID, sessionID)
	}
	return sessionID, nil
}

func GetApplicantByEmail(repoApplicant applicantRepo.ApplicantRepository, email string) (*dto.ApplicantOutput, error) {
	user, err := repoApplicant.GetByEmail(email)
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
	}, err
}

func GetEmployerByEmail(repoEmployer employerRepo.EmployerRepository, email string) (*dto.EmployerOutput, error) {
	user, err := repoEmployer.GetByEmail(email)
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
	}, err
}

func GetApplicantByID(repoApplicant applicantRepo.ApplicantRepository, id uint64) (*dto.ApplicantOutput, error) {
	user, err := repoApplicant.GetByID(id)
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
	}, err
}

func GetEmployerByID(repoEmployer employerRepo.EmployerRepository, id uint64) (*dto.EmployerOutput, error) {
	user, err := repoEmployer.GetByID(id)
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
	}, err
}
