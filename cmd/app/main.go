// Package main starts server and all handlers
package main

import (
	"database/sql"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	applicant_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	cvRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/repository"
	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	portfolioRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacanciesUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"

	// "github.com/go-park-mail-ru/2024_2_VKatuny/internal"

	//vacancies_repostory "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"

	_ "github.com/jackc/pgx/v5/stdlib"

	vacancies_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	//worker_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker/delivery"
	//worker_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker/repository"
)

func GetDBConnection() (*sql.DB, error) { //conf DatabaseConfig) (*sql.DB, error) {
	dbURL := "user=postgres dbname=postgres password=passIMO host=127.0.0.1 port=5432 sslmode=disable"
	// dbURL := fmt.Sprintf(
	// 	"postgres://%s:%s@%s:%d/%s?application_name=%s&search_path=%s&connect_timeout=%d",
	// 	conf.User,
	// 	conf.Password,
	// 	conf.Host,
	// 	conf.Port,
	// 	conf.DBName,
	// 	conf.AppName,
	// 	conf.Schema,
	// 	conf.ConnectionTimeout,
	// )

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}

// @title   uArt's API
// @version 1.0

// @contact.name Ifelsik
// @contact.url  https://github.com/Ifelsik

// @host     127.0.0.1:8080
// @BasePath /api/v1
func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	// Try gorilla?
	dbConnection, err := GetDBConnection() //(*config.Database)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	Mux := http.NewServeMux()

	//applicantRepository := applicant_repository.NewRepo()
	applicantRepository := applicant_repository.NewApplicantStorage(dbConnection)
	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewSessionStorage(dbConnection) // just do it!

	// Now your applicantRepository doesn't implement the IApplicantRepository interface
	// Oleg you should implement method Update for applicant_repository.PostgreSQLApplicantStorage
	applicantHandler := applicant_delivery.CreateApplicantHandler(applicantRepository, sessionApplicantRepository, conf.Server.GetAddress())
	Mux.Handle("/api/v1/registration/applicant", applicantHandler)

	//employerRepository := employer_repository.NewRepo()
	employerRepository := employer_repository.NewEmployerStorage(dbConnection)

	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository, sessionEmployerRepository, conf.Server.GetAddress())
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	//sessionApplicantRepository, sessionEmployerRepository := session_repository.NewRepo() // just do it!

	loginHandler := session_delivery.LoginHandler(
		sessionApplicantRepository,
		sessionEmployerRepository,
		applicantRepository,
		employerRepository,
		conf.Server.GetAddress(),
	)
	Mux.Handle("/api/v1/login", loginHandler)

	logoutHandler := session_delivery.LogoutHandler(sessionApplicantRepository,
		sessionEmployerRepository,
		applicantRepository,
		employerRepository)
	Mux.Handle("/api/v1/logout", logoutHandler)

	authorizedHandler := session_delivery.AuthorizedHandler(sessionApplicantRepository,
		sessionEmployerRepository,
		applicantRepository,
		employerRepository)
	Mux.Handle("/api/v1/authorized", authorizedHandler)

	// TODO: should be from db
	vacanciesRepository := vacancies_repository.NewVacanciesStorage(dbConnection)
	vacanciesListHandler := vacancies_delivery.GetVacanciesHandler(vacanciesRepository) //(&db.Vacancies)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

	repositories := &internal.Repositories{
		ApplicantRepository: applicant_repository.NewApplicantStorage(dbConnection), // implement IApplicantRepository. Add method `Update`
		PortfolioRepository: portfolioRepository.NewPortfolioStorage(dbConnection),  // implement IPortfolioRepository
		CVRepository:        cvRepository.NewCVStorage(dbConnection),                // also need this method
		VacanciesRepository: vacanciesRepository, 
	}
	usecases := &internal.Usecases{
		ApplicantUsecase: applicantUsecase.NewApplicantUsecase(logger, repositories),
		PortfolioUsecase: portfolioUsecase.NewPortfolioUsecase(logger, repositories),
		CVUsecase:        cvUsecase.NewCVsUsecase(logger, repositories),
		VacanciesUsecase: vacanciesUsecase.NewVacanciesUsecase(logger, repositories),
	}

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
	Mux.HandleFunc("/api/v1/employer/portfolio/", employerProfileHandlers.GetEmployerVacanciesHandler)

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
