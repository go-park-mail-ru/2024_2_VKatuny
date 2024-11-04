package internal

import (
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	// employerRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	applicantRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	// sessionRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	// vacanciesRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	portfolioRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	cvRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/repository"
)

type Repositories struct {
	// EmployerRepository employerRepository.EmployerRepository
	ApplicantRepository applicantRepository.IApplicantRepository
	// SessionRepository sessionRepository.SessionRepository
	// VacanciesRepository vacanciesRepository.VacanciesRepository
	PortfolioRepository portfolioRepository.IPortfolioRepository
	CVRepository        cvRepository.ICVsRepository
}

type Usecases struct {
	// employerUsecase employerUsecase.IEmployerUsecase
	ApplicantUsecase applicantUsecase.IApplicantUsecase
	PortfolioUsecase portfolioUsecase.IPortfolioUsecase
	CVUsecase        cvUsecase.ICVsUsecase
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }

