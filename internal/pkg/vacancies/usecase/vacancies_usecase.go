// Package usecase contains usecase for vacancies
package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/sirupsen/logrus"
)

var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

type VacanciesUsecase struct {
	logger              *logrus.Entry
	vacanciesRepository vacancies.IVacanciesRepository
	
}

func NewVacanciesUsecase(logger *logrus.Logger, repositories *internal.Repositories) *VacanciesUsecase {
	return &VacanciesUsecase{
		logger:              logrus.NewEntry(logger),
		vacanciesRepository: repositories.VacanciesRepository,
	}
}

func (u *VacanciesUsecase) ValidateQueryParameters(ctx context.Context, offsetStr, numStr string) (uint64, uint64, error) {
	fn := "VacanciesUsecase.ValidateQueryParameters"
	u.logger = utils.SetLoggerRequestID(ctx, u.logger)

	var err error
	offset, err1 := strconv.Atoi(offsetStr)

	if err1 != nil {
		u.logger.Errorf("%s: query parameter offset isn't a number: %s", fn, err1)
		offset = 0
		err = ErrOffsetIsNotANumber
	}

	num, err2 := strconv.Atoi(numStr)
	if err2 != nil {
		u.logger.Errorf("%s:query parameter num isn't a number: %s", fn, err2)
		num = 0
		err = ErrNumIsNotANumber // previous err will be overwritten
	}
	return uint64(offset), uint64(num), err
}

const (
	defaultVacanciesOffset = 0
	defaultVacanciesNum    = 10
)

func (vu *VacanciesUsecase) SearchVacancies(ctx context.Context, offsetStr, numStr, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error) {
	fn := "VacanciesUsecase.GetVacanciesWithOffset"
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)
	
	offset, num, err := vu.ValidateQueryParameters(ctx, offsetStr, numStr)
	if errors.Is(ErrOffsetIsNotANumber, err) {
		offset = defaultVacanciesOffset
	}
	if errors.Is(ErrNumIsNotANumber, err) {
		num = defaultVacanciesNum
	}
	if searchStr == "" && searchBy != "" {
		searchBy = ""
	}
	var vacancies []*dto.JSONVacancy
	vacancies, err = vu.vacanciesRepository.SearchAll(ctx, offset, offset+num, searchStr, group, searchBy)
	for _, vacancy := range vacancies {
		utils.EscapeHTMLStruct(vacancy)
	}
	if err != nil {
		return nil, err
	}
	vu.logger.Debugf("%s: got %d vacancies", fn, len(vacancies))
	return vacancies, nil
}

func (vu *VacanciesUsecase) GetVacanciesByEmployerID(ctx context.Context, employerID uint64) ([]*dto.JSONGetEmployerVacancy, error) {
	fn := "VacanciesUsecase.GetVacanciesByEmployerID"
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vacanciesModels, err := vu.vacanciesRepository.GetVacanciesByEmployerID(ctx, employerID)
	if err != nil {
		vu.logger.Errorf("function %s: unable to get vacancies: %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}

	vacancies := make([]*dto.JSONGetEmployerVacancy, 0, len(vacanciesModels))
	for _, vacancyModel := range vacanciesModels {
		utils.EscapeHTMLStruct(vacancyModel)
		vacancies = append(vacancies, &dto.JSONGetEmployerVacancy{
			ID:                   vacancyModel.ID,
			EmployerID:           vacancyModel.EmployerID,
			Salary:               vacancyModel.Salary,
			Position:             vacancyModel.Position,
			Description:          vacancyModel.Description,
			Location:             vacancyModel.Location,
			WorkType:             vacancyModel.WorkType,
			Avatar:               vacancyModel.Avatar,
			PositionCategoryName: vacancyModel.PositionCategoryName,
			CreatedAt:            vacancyModel.CreatedAt,
			UpdatedAt:            vacancyModel.UpdatedAt,
			CompressedAvatar:     vacancyModel.CompressedAvatar,
		})
	}
	return vacancies, nil
}

func (vu *VacanciesUsecase) CreateVacancy(ctx context.Context, vacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error) {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got creation request")
	vacancy.EmployerID = currentUser.ID
	utils.EscapeHTMLStruct(vacancy)
	createdVacancyID, err := vu.vacanciesRepository.Create(ctx, vacancy)
	if err != nil {
		vu.logger.Errorf("while creating in db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("vacancy created successfully with id %d", createdVacancyID)

	updatedVacancy, err := vu.vacanciesRepository.GetByID(ctx, createdVacancyID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	utils.EscapeHTMLStruct(updatedVacancy)
	vu.logger.Debugf("got updated vacancy with id %d", createdVacancyID)
	return updatedVacancy, nil
}

func (vu *VacanciesUsecase) GetVacancy(ctx context.Context, ID uint64) (*dto.JSONVacancy, error) {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vacancy, err := vu.vacanciesRepository.GetByID(ctx, ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, err
	}
	utils.EscapeHTMLStruct(vacancy)
	return vacancy, nil
}

func (vu *VacanciesUsecase) UpdateVacancy(ctx context.Context, ID uint64, vacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error) {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got update request")
	oldVacancy, err := vu.vacanciesRepository.GetByID(ctx, ID)
	utils.EscapeHTMLStruct(oldVacancy)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	vacancy.EmployerID = oldVacancy.EmployerID
	if vacancy.EmployerID != currentUser.ID {
		vu.logger.Debugf("not an owner tried to update vacancy, got %d expected %d", currentUser.ID, ID)
		return nil, fmt.Errorf(dto.MsgAccessDenied)
	}
	utils.EscapeHTMLStruct(vacancy)

	updatedVacancy, err := vu.vacanciesRepository.Update(ctx, ID, vacancy)
	if err != nil {
		vu.logger.Errorf("while updating in db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	utils.EscapeHTMLStruct(updatedVacancy)
	vu.logger.Debugf("successfully updated vacancy, got %d", updatedVacancy.ID)
	return updatedVacancy, nil
}

func (vu *VacanciesUsecase) DeleteVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vu.logger.WithFields(logrus.Fields{"employer_id": currentUser.ID, "user_type": currentUser.UserType}).Debug("got delete request")
	vacancy, err := vu.vacanciesRepository.GetByID(ctx, ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	utils.EscapeHTMLStruct(vacancy)
	if vacancy.EmployerID != currentUser.ID {
		vu.logger.Debugf("not an owner tried to delete vacancy, got %d expected %d", currentUser.ID, ID)
		return fmt.Errorf(dto.MsgAccessDenied)
	}
	err = vu.vacanciesRepository.Delete(ctx, ID)
	if err != nil {
		vu.logger.Errorf("while deleting from db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("successfully deleted vacancy with %d", ID)
	return nil
}

func (vu *VacanciesUsecase) SubscribeOnVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	if currentUser == nil {
		vu.logger.Errorf("user is not provided")
		return fmt.Errorf(dto.MsgUnauthorized)
	}

	err := vu.vacanciesRepository.Subscribe(ctx, ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while subscribing on db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("successfully subscribed on vacancy with ID: %d", ID)
	return nil
}

func (vu *VacanciesUsecase) UnsubscribeFromVacancy(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	if currentUser == nil {
		vu.logger.Errorf("user is not provided, currentUser = %v", currentUser)
		return fmt.Errorf(dto.MsgUnauthorized)
	}

	err := vu.vacanciesRepository.Unsubscribe(ctx, ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while unsubscribing on db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("successfully unsubscribed from vacancy with ID: %d", ID)
	return nil
}

func (vu *VacanciesUsecase) GetSubscriptionInfo(ctx context.Context, ID, applicantID uint64) (*dto.JSONVacancySubscriptionStatus, error) {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	isApplicantSubscribed, err := vu.vacanciesRepository.GetSubscriptionStatus(ctx, ID, applicantID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("got subscription status: %v, for vacancy %d and user %d", isApplicantSubscribed, ID, applicantID)
	return &dto.JSONVacancySubscriptionStatus{
		ID:           ID,
		ApplicantID:  applicantID,
		IsSubscribed: isApplicantSubscribed,
	}, nil
}

func (vu *VacanciesUsecase) GetVacancySubscribers(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) (*dto.JSONVacancySubscribers, error) {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	if currentUser == nil {
		vu.logger.Errorf("user is not provided, currentUser = %v", currentUser)
		return nil, fmt.Errorf(dto.MsgUnauthorized)
	}

	vacancy, err := vu.vacanciesRepository.GetByID(ctx, ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}

	if currentUser.UserType != dto.UserTypeEmployer || currentUser.ID != vacancy.EmployerID {
		vu.logger.Errorf("user is not employer, currentUser = %v", currentUser)
		return nil, fmt.Errorf(dto.MsgUnauthorized)
	}

	subscribersModel, err := vu.vacanciesRepository.GetSubscribersList(ctx, ID)
	if err != nil {
		vu.logger.Errorf("while getting from db got err %s", err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}

	subscribers := make([]*dto.JSONGetApplicantProfile, 0, len(subscribersModel))
	for _, subscriberModel := range subscribersModel {
		utils.EscapeHTMLStruct(subscriberModel)
		subscribers = append(subscribers, &dto.JSONGetApplicantProfile{
			ID:               subscriberModel.ID,
			FirstName:        subscriberModel.FirstName,
			LastName:         subscriberModel.LastName,
			City:             subscriberModel.CityName,
			BirthDate:        subscriberModel.BirthDate,
			Avatar:           subscriberModel.PathToProfileAvatar,
			Contacts:         subscriberModel.Contacts,
			Education:        subscriberModel.Education,
			CompressedAvatar: subscriberModel.CompressedAvatar,
		})
	}

	vu.logger.Debugf("successfully got %d subscribers for vacancy %d", len(subscribers), ID)
	return &dto.JSONVacancySubscribers{
		ID:          ID,
		Subscribers: subscribers,
	}, nil
}

func (vu *VacanciesUsecase) GetApplicantFavoriteVacancies(ctx context.Context, applicantID uint64) ([]*dto.JSONGetEmployerVacancy, error) {
	fn := "VacanciesUsecase.GetApplicantFavoriteVacancies"
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	vacanciesModels, err := vu.vacanciesRepository.GetApplicantFavoriteVacancies(ctx, applicantID)
	if err != nil {
		vu.logger.Errorf("function %s: unable to get vacancies: %s", fn, err)
		return nil, fmt.Errorf(dto.MsgDataBaseError)
	}

	vacancies := make([]*dto.JSONGetEmployerVacancy, 0, len(vacanciesModels))
	for _, vacancyModel := range vacanciesModels {
		vacancies = append(vacancies, &dto.JSONGetEmployerVacancy{
			ID:                   vacancyModel.ID,
			EmployerID:           vacancyModel.EmployerID,
			Salary:               vacancyModel.Salary,
			Position:             vacancyModel.Position,
			Description:          vacancyModel.Description,
			WorkType:             vacancyModel.WorkType,
			Avatar:               vacancyModel.Avatar,
			PositionCategoryName: vacancyModel.PositionCategoryName,
			CreatedAt:            vacancyModel.CreatedAt,
			UpdatedAt:            vacancyModel.UpdatedAt,
			CompressedAvatar:     vacancyModel.CompressedAvatar,
		})
	}
	return vacancies, nil
}

func (vu *VacanciesUsecase) AddIntoFavorite(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)

	if currentUser == nil {
		vu.logger.Errorf("user is not provided")
		return fmt.Errorf(dto.MsgUnauthorized)
	}

	err := vu.vacanciesRepository.MakeFavorite(ctx, ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while adding into favorite on db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("successfully adding into favorite on vacancy with ID: %d", ID)
	return nil
}

func (vu *VacanciesUsecase) Unfavorite(ctx context.Context, ID uint64, currentUser *dto.UserFromSession) error {
	vu.logger = utils.SetLoggerRequestID(ctx, vu.logger)
	if currentUser == nil {
		vu.logger.Errorf("user is not provided")
		return fmt.Errorf(dto.MsgUnauthorized)
	}

	err := vu.vacanciesRepository.Unfavorite(ctx, ID, currentUser.ID)
	if err != nil {
		vu.logger.Errorf("while adding into favorite on db got err %s", err)
		return fmt.Errorf(dto.MsgDataBaseError)
	}
	vu.logger.Debugf("successfully adding into favorite on vacancy with ID: %d", ID)
	return nil
}
