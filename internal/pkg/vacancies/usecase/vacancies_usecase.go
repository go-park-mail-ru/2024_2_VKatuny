// Package usecase contains usecase for vacancies
package usecase

import (
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/sirupsen/logrus"
)

var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

type VacanciesUsecase struct {
	logger              *logrus.Logger
	vacanciesRepository vacancies.IVacanciesRepository
}

func NewVacanciesUsecase(logger *logrus.Logger, repositories *internal.Repositories) *VacanciesUsecase {
	return &VacanciesUsecase{
		logger:              logger,
		vacanciesRepository: repositories.VacanciesRepository,
	}
}

func ValidateRequestParams(offsetStr, numStr string) (uint64, uint64, error) {
	var err error
	offset, err1 := strconv.Atoi(offsetStr)

	if err1 != nil {
		offset = 0
		err = ErrOffsetIsNotANumber
	}

	num, err2 := strconv.Atoi(numStr)
	if err2 != nil {
		num = 0
		err = ErrNumIsNotANumber // previous err will be overwritten
	}
	return uint64(offset), uint64(num), err
}

func (vu *VacanciesUsecase) GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONGetEmployerVacancy, error) {
	fn := "VacanciesUsecase.GetVacanciesByEmployerID"
	vacanciesModels, err := vu.vacanciesRepository.GetVacanciesByEmployerID(employerID)
	if err != nil {
		vu.logger.Errorf("function %s: unable to get vacancies: %s", fn, err)
		return nil, err
	}

	vacancies := make([]*dto.JSONGetEmployerVacancy, 0, len(vacanciesModels))
	for _, vacancyModel := range vacanciesModels {
		vacancies = append(vacancies, &dto.JSONGetEmployerVacancy{
			ID:          vacancyModel.ID,
			EmployerID:  vacancyModel.EmployerID,
			Salary:      vacancyModel.Salary,
			Position:    vacancyModel.Position,
			Description: vacancyModel.Description,
			WorkType:    vacancyModel.WorkType,
			Avatar:      vacancyModel.Avatar,
			CreatedAt:   vacancyModel.CreatedAt,
		})
	}
	return vacancies, nil
}

func (vu *VacanciesUsecase) CreateVacancy(vacancy *dto.JSONVacancy, currentUser *dto.SessionUser) (*dto.JSONVacancy, error) {
	// TODO: need to validate vacancy && currentUser is not nil

	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got creation request")
	vacancy.EmployerID = currentUser.ID
	createdVacancyID, err := vu.vacanciesRepository.Create(vacancy)
	if err != nil {
		return nil, err
	}
	fmt.Println(createdVacancyID, "123213")
	return vu.vacanciesRepository.GetByID(createdVacancyID)
}

func (vu *VacanciesUsecase) GetVacancy(ID uint64) (*dto.JSONVacancy, error) {
	vacancy, err := vu.vacanciesRepository.GetByID(ID)
	if err != nil {
		return nil, err
	}
	return vacancy, nil
}

func (vu *VacanciesUsecase) UpdateVacancy(ID uint64, vacancy *dto.JSONVacancy, currentUser *dto.SessionUser) (*dto.JSONVacancy, error) {
	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got update request")
	oldVacancy, err := vu.vacanciesRepository.GetByID(ID)
	fmt.Println(oldVacancy)
	if err != nil {
		return nil, err
	}
	vacancy.EmployerID = oldVacancy.EmployerID
	if vacancy.EmployerID != currentUser.ID {
		vu.logger.Debugf("not an owner tried to update vacancy, got %d expected %d", currentUser.ID, ID)
		return nil, commonerrors.ErrUnauthorized
	}

	updatedVacancy, err := vu.vacanciesRepository.Update(ID, vacancy)
	if err != nil {
		return nil, err
	}
	vu.logger.Debugf("successfully updated vacancy, got %d", updatedVacancy.ID)
	return updatedVacancy, nil
}

func (vu *VacanciesUsecase) DeleteVacancy(ID uint64, currentUser *dto.SessionUser) error {
	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got delete request")
	vacancy, err := vu.vacanciesRepository.GetByID(ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return err
	}
	if vacancy.EmployerID != currentUser.ID {
		vu.logger.Debugf("not an owner tried to delete vacancy, got %d expected %d", currentUser.ID, ID)
		return commonerrors.ErrUnauthorized
	}
	err = vu.vacanciesRepository.Delete(ID)
	if err != nil {
		vu.logger.Errorf("while deleting from db got err %s", err)
		return err
	}
	return nil
}

func (vu *VacanciesUsecase) SubscribeOnVacancy(ID uint64, currentUser *dto.SessionUser) error {
	if currentUser == nil {
		vu.logger.Errorf("user is not provided")
		return commonerrors.ErrUnauthorized
	}

	err := vu.vacanciesRepository.Subscribe(ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while subscribing on db got err %s", err)
		return err
	}
	return nil
}

func (vu *VacanciesUsecase) UnsubscribeFromVacancy(ID uint64, currentUser *dto.SessionUser) error {
	if currentUser == nil {
		vu.logger.Errorf("user is not provided, currentUser = %v", currentUser)
		return commonerrors.ErrUnauthorized
	}

	err := vu.vacanciesRepository.Unsubscribe(ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while unsubscribing on db got err %s", err)
		return err
	}
	return nil
}

func (vu *VacanciesUsecase) GetSubscriptionInfo(ID, applicantID uint64) (*dto.JSONVacancySubscriptionStatus, error) {
	isApplicantSubscribed, err := vu.vacanciesRepository.GetSubscriptionStatus(ID, applicantID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, err
	}
	return &dto.JSONVacancySubscriptionStatus{
		ID:           ID,
		ApplicantID:  applicantID,
		IsSubscribed: isApplicantSubscribed,
	}, nil
}

func (vu *VacanciesUsecase) GetVacancySubscribers(ID uint64, currentUser *dto.SessionUser) (*dto.JSONVacancySubscribers, error) {
	if currentUser == nil {
		vu.logger.Errorf("user is not provided, currentUser = %v", currentUser)
		return nil, commonerrors.ErrUnauthorized
	}

	vacancy, err := vu.vacanciesRepository.GetByID(ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, err
	}

	if currentUser.UserType != dto.UserTypeEmployer || currentUser.ID != vacancy.EmployerID {
		vu.logger.Errorf("user is not applicant, currentUser = %v", currentUser)
		return nil, commonerrors.ErrUnauthorized
	}

	subscribersModel, err := vu.vacanciesRepository.GetSubscribersList(ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, err
	}

	subscribers := make([]*dto.JSONGetApplicantProfile, 0, len(subscribersModel))
	for _, subscriberModel := range subscribersModel {
		subscribers = append(subscribers, &dto.JSONGetApplicantProfile{
			ID:        subscriberModel.ID,
			FirstName: subscriberModel.FirstName,
			LastName:  subscriberModel.LastName,
			City:      subscriberModel.CityName,
			BirthDate: subscriberModel.BirthDate,
			Avatar:    subscriberModel.PathToProfileAvatar,
			Contacts:  subscriberModel.Contacts,
			Education: subscriberModel.Education,
		})
	}
	return &dto.JSONVacancySubscribers{
		ID:          ID,
		Subscribers: subscribers,
	}, nil
}
