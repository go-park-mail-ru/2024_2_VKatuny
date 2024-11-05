package internal

type Repositories struct {
	EmployerRepository  interface{}
	ApplicantRepository interface{}
	// SessionRepository sessionRepository.SessionRepository
	PortfolioRepository interface{}
	CVRepository        interface{}
	VacanciesRepository interface{}
}

type Usecases struct {
	EmployerUsecase  interface{}
	ApplicantUsecase interface{}
	PortfolioUsecase interface{}
	CVUsecase        interface{}
	VacanciesUsecase interface{}
}

// type Handlers struct {
// 	ApplicantProfileHandlers *ApplicantProfileHandlers
// }

