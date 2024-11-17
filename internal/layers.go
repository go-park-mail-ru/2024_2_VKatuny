package internal

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger       *logrus.Logger
	Repositories *Repositories
	Usecases     *Usecases
}

type Repositories struct {
	EmployerRepository         interface{}
	ApplicantRepository        interface{}
	PortfolioRepository        interface{}
	CVRepository               cvs.ICVsRepository
	VacanciesRepository        vacancies.IVacanciesRepository
	SessionApplicantRepository session.ISessionRepository
	SessionEmployerRepository  session.ISessionRepository
	FileLoadingRepository      fileloading.IFileLoadingRepository
}

type Usecases struct {
	EmployerUsecase    interface{}
	ApplicantUsecase   interface{}
	PortfolioUsecase   interface{}
	CVUsecase          cvs.ICVsUsecase
	VacanciesUsecase   vacancies.IVacanciesUsecase
	FileLoadingUsecase fileloading.IFileLoadingUsecase
	SessionUsecase     session.ISessionUsecase
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }
