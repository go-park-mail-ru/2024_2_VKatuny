// Package usecase contains usecase for session
package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

var ErrEmptyCookie = fmt.Errorf("client have an empty cookie")

type sessionUsecase struct {
	logger               *logrus.Entry
	applicantRepo        applicant.IApplicantRepository
	employerRepo         employer.IEmployerRepository
	sessionRepoApplicant session.ISessionRepository
	sessionRepoEmployer  session.ISessionRepository
}

func NewSessionUsecase(logger *logrus.Logger, repositories *internal.Repositories) *sessionUsecase {
	return &sessionUsecase{
		logger:               &logrus.Entry{Logger: logger},
		sessionRepoApplicant: repositories.SessionApplicantRepository,
		sessionRepoEmployer:  repositories.SessionEmployerRepository,
	}
}

// TODO: should return user
func (u *sessionUsecase) Login(user *dto.JSONLoginForm) (*dto.UserWithSession, error) {
	fn := "sessionUsecase.Login"
	if user == nil {
		u.logger.Errorf("%s: user is nil", fn)
		return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
	}
	u.logger.Debugf("%s: logging in as %v", fn, user)
	switch user.UserType {
	case dto.UserTypeApplicant:
		applicant, err := u.applicantRepo.GetByEmail(user.Email)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		if !utils.EqualHashedPasswords(applicant.PasswordHash, user.Password) {
			u.logger.Errorf("%s: password comparison failed", fn)
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}

		sessionID := utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeApplicant)
		err = u.sessionRepoApplicant.Create(applicant.ID, sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}

		// TODO: logger with method from other PR
		u.logger.
			WithFields(logrus.Fields{"user_id": applicant.ID, "user_type": user.UserType}).
			Debugf("%s: successfully logged in", fn)

		return &dto.UserWithSession{
			ID:        applicant.ID,
			UserType:  dto.UserTypeApplicant,
			SessionID: sessionID,
		}, nil
	case dto.UserTypeEmployer:
		employer, err := u.employerRepo.GetByEmail(user.Email)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}
		if !utils.EqualHashedPasswords(employer.PasswordHash, user.Password) {
			u.logger.Errorf("%s: password comparison failed", fn)
			return nil, fmt.Errorf(dto.MsgWrongLoginOrPassword)
		}

		sessionID := utils.GenerateSessionToken(utils.TokenLength, dto.UserTypeEmployer)

		err = u.sessionRepoEmployer.Create(employer.ID, sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}

		u.logger.
			WithFields(logrus.Fields{"user_id": employer.ID, "user_type": user.UserType}).
			Debugf("%s: successfully logged in", fn)

		return &dto.UserWithSession{
			ID:        employer.ID,
			UserType:  dto.UserTypeEmployer,
			SessionID: sessionID,
		}, nil
	}
	u.logger.Errorf("%s: bad user type", fn)
	return nil, fmt.Errorf(dto.MsgBadUserType)
}

func (u *sessionUsecase) Logout(userType string, sessionID string) (*dto.JSONUser, error) {
	fn := "sessionUsecase.Logout"
	switch userType {
	case dto.UserTypeApplicant:
		applicantID, err := u.sessionRepoApplicant.GetUserIdBySession(sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}
		err = u.sessionRepoApplicant.Delete(sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}
		u.logger.
			WithFields(logrus.Fields{"user_id": applicantID, "user_type": userType}).
			Debugf("%s: successfully logged out", fn)
		return &dto.JSONUser{
			ID:       applicantID,
			UserType: dto.UserTypeApplicant,
		}, nil
	case dto.UserTypeEmployer:
		employerID, err := u.sessionRepoEmployer.GetUserIdBySession(sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}
		err = u.sessionRepoEmployer.Delete(sessionID)
		if err != nil {
			u.logger.Errorf("%s: got err %s", fn, err)
			return nil, fmt.Errorf(dto.MsgDataBaseError)
		}
		u.logger.
			WithFields(logrus.Fields{"user_id": employerID, "user_type": userType}).
			Debugf("%s: successfully logged out", fn)
		return &dto.JSONUser{
			ID:       employerID,
			UserType: dto.UserTypeEmployer,
		}, nil
	}
	return nil, fmt.Errorf(dto.MsgBadUserType)
}

func (u *sessionUsecase) CheckAuthorization(userType, sessionID string) (uint64, error) {
	if sessionID == "" {
		return 0, fmt.Errorf(dto.MsgBadCookie)
	}
	if userType == dto.UserTypeApplicant {
		applicantID, err := u.sessionRepoApplicant.GetUserIdBySession(sessionID)
		if err != nil {
			return 0, fmt.Errorf(dto.MsgNoUserWithSession)
		}
		return applicantID, nil
	} else if userType == dto.UserTypeEmployer {
		employerID, err := u.sessionRepoEmployer.GetUserIdBySession(sessionID)
		if err != nil {
			return 0, fmt.Errorf(dto.MsgNoUserWithSession)
		}
		return employerID, nil
	}
	return 0, fmt.Errorf(dto.MsgBadUserType)
}

// func GetApplicantByEmail(repoApplicant applicantRepo.IApplicantRepository, email string) (*dto.ApplicantOutput, error) {
// 	user, err := repoApplicant.GetByEmail(email)
// 	if err != nil {
// 		return nil, fmt.Errorf(dto.MsgDataBaseError)
// 	}
// 	return &dto.ApplicantOutput{
// 		ID:                  user.ID,
// 		FirstName:           user.FirstName,
// 		LastName:            user.LastName,
// 		CityName:            user.CityName,
// 		BirthDate:           user.BirthDate,
// 		PathToProfileAvatar: user.PathToProfileAvatar,
// 		Contacts:            user.Contacts,
// 		Education:           user.Education,
// 		Email:               user.Email,
// 		CreatedAt:           user.CreatedAt,
// 		UpdatedAt:           user.UpdatedAt,
// 	}, nil
// }

// func GetEmployerByEmail(repoEmployer employerRepo.EmployerRepository, email string) (*dto.EmployerOutput, error) {
// 	user, err := repoEmployer.GetByEmail(email)
// 	if err != nil {
// 		return nil, fmt.Errorf(dto.MsgDataBaseError)
// 	}
// 	return &dto.EmployerOutput{
// 		ID:                  user.ID,
// 		FirstName:           user.FirstName,
// 		LastName:            user.LastName,
// 		CityName:            user.CityName,
// 		Position:            user.Position,
// 		CompanyName:         user.CompanyName,
// 		CompanyDescription:  user.CompanyDescription,
// 		CompanyWebsite:      user.CompanyWebsite,
// 		PathToProfileAvatar: user.PathToProfileAvatar,
// 		Contacts:            user.Contacts,
// 		Email:               user.Email,
// 		PasswordHash:        user.PasswordHash,
// 		CreatedAt:           user.CreatedAt,
// 		UpdatedAt:           user.UpdatedAt,
// 	}, nil
// }

// func GetApplicantByID(repoApplicant applicantRepo.IApplicantRepository, id uint64) (*dto.ApplicantOutput, error) {
// 	user, err := repoApplicant.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf(dto.MsgDataBaseError)
// 	}
// 	return &dto.ApplicantOutput{
// 		ID:                  user.ID,
// 		FirstName:           user.FirstName,
// 		LastName:            user.LastName,
// 		CityName:            user.CityName,
// 		BirthDate:           user.BirthDate,
// 		PathToProfileAvatar: user.PathToProfileAvatar,
// 		Contacts:            user.Contacts,
// 		Education:           user.Education,
// 		Email:               user.Email,
// 		CreatedAt:           user.CreatedAt,
// 		UpdatedAt:           user.UpdatedAt,
// 	}, nil
// }

// func GetEmployerByID(repoEmployer employerRepo.EmployerRepository, id uint64) (*dto.EmployerOutput, error) {
// 	user, err := repoEmployer.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf(dto.MsgDataBaseError)
// 	}
// 	return &dto.EmployerOutput{
// 		ID:                  user.ID,
// 		FirstName:           user.FirstName,
// 		LastName:            user.LastName,
// 		CityName:            user.CityName,
// 		Position:            user.Position,
// 		CompanyName:         user.CompanyName,
// 		CompanyDescription:  user.CompanyDescription,
// 		CompanyWebsite:      user.CompanyWebsite,
// 		PathToProfileAvatar: user.PathToProfileAvatar,
// 		Contacts:            user.Contacts,
// 		Email:               user.Email,
// 		PasswordHash:        user.PasswordHash,
// 		CreatedAt:           user.CreatedAt,
// 		UpdatedAt:           user.UpdatedAt,
// 	}, nil
// }
