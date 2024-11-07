// Package usecase contains usecase for session
package usecase

import (
	"fmt"
	"net/http"

	applicantRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	employerRepo "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
)

var ErrEmptyCookie = fmt.Errorf("client have an empty cookie")

func CheckAuthorization(newUserInput *dto.JSONLogoutForm, session *http.Cookie, sessionRepoApplicant session.ISessionRepository, sessionRepoEmployer session.ISessionRepository) (uint64, error) {
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
	return 0, fmt.Errorf(dto.MsgBadUserType)
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
func LoginValidate(newUserInput *dto.JSONLoginForm, repoApplicant applicantRepo.IApplicantRepository, repoEmployer employerRepo.EmployerRepository) (user *dto.UserIDAndType, err error) {
	// TODO: add logging
	if newUserInput.UserType == dto.UserTypeApplicant {
		var worker *models.Applicant
		worker, err = repoApplicant.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		if !utils.EqualHashedPasswords(worker.PasswordHash, newUserInput.Password) {
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		user = &dto.UserIDAndType{ID: worker.ID, UserType: dto.UserTypeApplicant}
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		var employer *models.Employer
		employer, err = repoEmployer.GetByEmail(newUserInput.Email)
		if err != nil {
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		if !utils.EqualHashedPasswords(employer.PasswordHash, newUserInput.Password) {
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		user = &dto.UserIDAndType{ID: employer.ID, UserType: dto.UserTypeEmployer}
	}
	return user, err
}

// LogoutValidate tries to remove session from db
// TODO: rename function to more accurate meaning
func LogoutValidate(newUserInput *dto.JSONLogoutForm, session string, sessionRepoApplicant, sessionRepoEmployer session.ISessionRepository) (uint64, error) {
	if newUserInput.UserType == dto.UserTypeApplicant {
		id, err := sessionRepoApplicant.GetUserIdBySession(session)
		if err != nil {
			return 0, fmt.Errorf(dto.MsgDataBaseError)
		}
		err = sessionRepoApplicant.Delete(session)
		if err == nil {
			return id, nil
		}
	} else if newUserInput.UserType == dto.UserTypeEmployer {
		id, err := sessionRepoEmployer.GetUserIdBySession(session)
		if err != nil {
			return 0, fmt.Errorf(dto.MsgDataBaseError)
		}
		err = sessionRepoEmployer.Delete(session)
		if err == nil {
			return id, nil
		}
	}
	return 0, fmt.Errorf(dto.MsgBadUserType)
}

func AddSession(sessionRepoApplicant, sessionRepoEmployer session.ISessionRepository, user *dto.UserIDAndType) (string, error) {
	var sessionID string
	switch user.UserType {
	case dto.UserTypeApplicant:
		sessionID = utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeApplicant)
		if err := sessionRepoApplicant.Create(user.ID, sessionID); err != nil {
			return "", fmt.Errorf(dto.MsgDataBaseError)
		}
	case dto.UserTypeEmployer:
		sessionID = utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeEmployer)
		if err := sessionRepoEmployer.Create(user.ID, sessionID); err != nil {
			return "", fmt.Errorf(dto.MsgDataBaseError)
		}
	}
	return sessionID, nil
}

func GetApplicantByEmail(repoApplicant applicantRepo.IApplicantRepository, email string) (*dto.ApplicantOutput, error) {
	user, err := repoApplicant.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return &dto.ApplicantOutput{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		CityName:            user.CityName,
		BirthDate:           user.BirthDate,
		PathToProfileAvatar: user.PathToProfileAvatar,
		Contacts:            user.Contacts,
		Education:           user.Education,
		Email:               user.Email,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, nil
}

func GetEmployerByEmail(repoEmployer employerRepo.EmployerRepository, email string) (*dto.EmployerOutput, error) {
	user, err := repoEmployer.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
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
	}, nil
}

func GetApplicantByID(repoApplicant applicantRepo.IApplicantRepository, id uint64) (*dto.ApplicantOutput, error) {
	user, err := repoApplicant.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	return &dto.ApplicantOutput{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		CityName:            user.CityName,
		BirthDate:           user.BirthDate,
		PathToProfileAvatar: user.PathToProfileAvatar,
		Contacts:            user.Contacts,
		Education:           user.Education,
		Email:               user.Email,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, nil
}

func GetEmployerByID(repoEmployer employerRepo.EmployerRepository, id uint64) (*dto.EmployerOutput, error) {
	user, err := repoEmployer.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
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
	}, nil
}
