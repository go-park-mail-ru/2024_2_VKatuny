// Package main starts server and all handlers
package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	applicant_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	cvDelivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	cvRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/repository"
	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	portfolioRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacanciesUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"

	vacancies_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
)

// @title   uArt's API
// @version 1.0

// @contact.name Ifelsik
// @contact.url  https://github.com/Ifelsik

// @host     127.0.0.1:8080
// @BasePath /api/v1
func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	dbConnection, err := utils.GetDBConnection(conf.DataBase.GetDSN()) 
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	Mux := http.NewServeMux()

	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewSessionStorage(dbConnection)
	repositories := &internal.Repositories{
		ApplicantRepository:        applicant_repository.NewApplicantStorage(dbConnection),
		PortfolioRepository:        portfolioRepository.NewPortfolioStorage(dbConnection),
		CVRepository:               cvRepository.NewCVStorage(dbConnection),
		VacanciesRepository:        vacancies_repository.NewVacanciesStorage(dbConnection),
		EmployerRepository:         employer_repository.NewEmployerStorage(dbConnection),
		SessionApplicantRepository: sessionApplicantRepository,
		SessionEmployerRepository:  sessionEmployerRepository,
	}
	usecases := &internal.Usecases{
		ApplicantUsecase: applicantUsecase.NewApplicantUsecase(logger, repositories),
		PortfolioUsecase: portfolioUsecase.NewPortfolioUsecase(logger, repositories),
		CVUsecase:        cvUsecase.NewCVsUsecase(logger, repositories),
		VacanciesUsecase: vacanciesUsecase.NewVacanciesUsecase(logger, repositories),
		EmployerUsecase:  employerUsecase.NewEmployerUsecase(logger, repositories),
	}
	app := &internal.App{
		Logger:       logger,
		Repositories: repositories,
		Usecases:     usecases,
	}


	applicantHandler := applicant_delivery.CreateApplicantHandler(repositories.ApplicantRepository, sessionApplicantRepository, conf.Server.GetAddress())
	Mux.Handle("/api/v1/registration/applicant", applicantHandler)

	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository, sessionEmployerRepository, conf.Server.GetAddress())
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	sessionHandlers := session_delivery.NewSessionHandlers(app)
	Mux.HandleFunc("/api/v1/login", sessionHandlers.Login)
	Mux.HandleFunc("/api/v1/logout", sessionHandlers.Logout)
	Mux.HandleFunc("/api/v1/authorized", sessionHandlers.IsAuthorized)

	vacanciesListHandler := vacancies_delivery.GetVacanciesHandler(vacanciesRepository)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)


	applicantProfileHandlers, err := applicant_delivery.NewApplicantProfileHandlers(logger, usecases)
	if err != nil {
		logger.Fatal(err)
	}
	Mux.HandleFunc("/api/v1/applicant/profile/", applicantProfileHandlers.ApplicantProfileHandler)
	Mux.HandleFunc("/api/v1/applicant/portfolio/", applicantProfileHandlers.GetApplicantPortfoliosHandler)
	Mux.HandleFunc("/api/v1/applicant/cv/", applicantProfileHandlers.GetApplicantCVsHandler)

	employerProfileHandlers, err := employer_delivery.NewEmployerProfileHandlers(logger, usecases)
	if err != nil {
		logger.Fatal(err)
	}
	Mux.HandleFunc("/api/v1/employer/profile/", employerProfileHandlers.EmployerProfileHandler)
	Mux.HandleFunc("/api/v1/employer/vacancies/", employerProfileHandlers.GetEmployerVacanciesHandler)

	cvsHandlers := cvDelivery.NewCVsHandler(app)
	Mux.HandleFunc("/api/v1/cv/", cvsHandlers.CVsRESTHandler)

	vacanciesHandlers := vacancies_delivery.NewVacanciesHandlers(app)
	Mux.HandleFunc("/api/v1/vacancy/", vacanciesHandlers.VacanciesRESTHandler)
	Mux.HandleFunc("/api/v1/vacancy/subscription/", vacanciesHandlers.VacanciesSubscribeRESTHandler)
	Mux.HandleFunc("/api/v1/vacancy/subscribers/", vacanciesHandlers.GetVacancySubscribersHandler)

	// Wrapped multiplexer
	// Mux implements http.Handler interface so it's possible to wrap
	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.GetFrontURI())
	handlers = middleware.Panic(handlers)
	handlers = middleware.AccessLogger(handlers, logger)
	handlers = middleware.SetContext(handlers, logger)
	logger.Infof("Server is starting at %s", conf.Server.GetAddress())
	err = http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
