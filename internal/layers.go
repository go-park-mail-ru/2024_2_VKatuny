package internal

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger         *logrus.Logger
	BackendAddress string
	Repositories   *Repositories
	Usecases       *Usecases
}

type Repositories struct {
	EmployerRepository         employer.IEmployerRepository
	ApplicantRepository        interface{}
	PortfolioRepository        interface{}
	CVRepository               cvs.ICVsRepository
	VacanciesRepository        vacancies.IVacanciesRepository
	SessionApplicantRepository session.ISessionRepository
	SessionEmployerRepository  session.ISessionRepository
}

type Usecases struct {
	EmployerUsecase  employer.IEmployerUsecase
	ApplicantUsecase interface{}
	PortfolioUsecase interface{}
	CVUsecase        cvs.ICVsUsecase
	VacanciesUsecase vacancies.IVacanciesUsecase
	SessionUsecase   session.ISessionUsecase
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }
