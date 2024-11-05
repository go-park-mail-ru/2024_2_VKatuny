// Package usecase contains usecase for vacancies
package usecase

import (
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	vacanciesRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/sirupsen/logrus"
)

var ErrOffsetIsNotANumber = fmt.Errorf("query parameter offset isn't a number")
var ErrNumIsNotANumber = fmt.Errorf("query parameter num isn't a number")

type IVacanciesUsecase interface {
	GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONGetEmployerVacancy, error)
}

type VacanciesUsecase struct {
	logger              *logrus.Logger
	vacanciesRepository vacanciesRepository.Repository
}

func NewVacanciesUsecase(logger *logrus.Logger, repositories *internal.Repositories) *VacanciesUsecase {
	vacanciesRepository, ok := repositories.VacanciesRepository.(vacanciesRepository.Repository)
	if !ok {
		return nil
	}
	return &VacanciesUsecase{
		logger:              logger,
		vacanciesRepository: vacanciesRepository,
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

	vacancies := make([]*dto.JSONGetEmployerVacancy, len(vacanciesModels))
	for _, vacancyModel := range vacanciesModels {
		vacancies = append(vacancies, &dto.JSONGetEmployerVacancy{
			ID:          vacancyModel.ID,
			EmployerID:  vacancyModel.EmployerID,
			Salary:      vacancyModel.Salary,
			Position:    vacancyModel.Position,
			Description: vacancyModel.Description,
			WorkType:    vacancyModel.WorkType,
			Avatar:      vacancyModel.Logo,
			CreatedAt:   vacancyModel.CreatedAt,
		})
	}
	return vacancies, nil
}
