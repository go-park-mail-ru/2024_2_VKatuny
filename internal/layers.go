package internal

import (
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer"
	fileloading "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies"
	authClient "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	notificationsmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"
	"github.com/sirupsen/logrus"
)

type App struct {
	Logger         *logrus.Logger
	BackendAddress string
	Repositories   *Repositories
	Usecases       *Usecases
	Microservices  *Microservices
}

type Repositories struct {
	EmployerRepository         employer.IEmployerRepository
	ApplicantRepository        applicant.IApplicantRepository
	PortfolioRepository        portfolio.IPortfolioRepository
	CVRepository               cvs.ICVsRepository
	VacanciesRepository        vacancies.IVacanciesRepository
	FileLoadingRepository      fileloading.IFileLoadingRepository
}

type Usecases struct {
	EmployerUsecase    employer.IEmployerUsecase
	ApplicantUsecase   applicant.IApplicantUsecase
	PortfolioUsecase   portfolio.IPortfolioUsecase
	CVUsecase          cvs.ICVsUsecase
	VacanciesUsecase   vacancies.IVacanciesUsecase
	FileLoadingUsecase fileloading.IFileLoadingUsecase
}

type Microservices struct {
	Auth     authClient.AuthorizationClient
	Compress compressmicroservice.CompressServiceClient
	Notifications notificationsmicroservice.NotificationsServiceClient
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }
