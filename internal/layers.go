package internal

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
)

type Repositories struct {
	EmployerRepository  interface{}
	ApplicantRepository interface{}
	// SessionRepository sessionRepository.SessionRepository
	PortfolioRepository interface{}
	CVRepository        cvs.ICVsRepository
	VacanciesRepository vacancies.IVacanciesRepository
}

type Usecases struct {
	EmployerUsecase  interface{}
	ApplicantUsecase interface{}
	PortfolioUsecase interface{}
	CVUsecase        cvs.ICVsUsecase
	VacanciesUsecase vacancies.IVacanciesUsecase
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }

