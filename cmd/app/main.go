// Package main starts server and all handlers
package main

import (
	"database/sql"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	applicant_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacancies_repostory "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	dbConnection, err := GetDBConnection() //(*config.Database)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	Mux := http.NewServeMux()

	applicantRepository := applicant_repository.NewRepo()
	applicantRepository1 := applicant_repository.NewApplicantStorage(dbConnection)

	applicantHandler := applicant_delivery.CreateApplicantHandler(applicantRepository1)
	Mux.Handle("/api/v1/registration/applicant", applicantHandler)

	employerRepository := employer_repository.NewRepo()
	employerRepository1 := employer_repository.NewEmployerStorage(dbConnection)

	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository1)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewRepo() // just do it!
	loginHandler := session_delivery.LoginHandler(
		sessionApplicantRepository,
		sessionEmployerRepository,
		applicantRepository,
		employerRepository,
		conf.Server.GetAddress(),
	)
	Mux.Handle("/api/v1/login", loginHandler)

	logoutHandler := session_delivery.LogoutHandler(sessionApplicantRepository, sessionEmployerRepository)
	Mux.Handle("/api/v1/logout", logoutHandler)

	authorizedHandler := session_delivery.AuthorizedHandler(sessionApplicantRepository, sessionEmployerRepository)
	Mux.Handle("/api/v1/authorized", authorizedHandler)

	vacanciesRepository := vacancies_repostory.NewRepo()
	vacanciesListHandler := vacancies_delivery.GetVacanciesHandler(vacanciesRepository) //(&db.Vacancies)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

	// Wrapped multiplexer
	// Mux implements http.Handler interface so it's possible to wrap
	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.GetHostWithScheme())
	handlers = middleware.Panic(handlers)
	handlers = middleware.AccessLogger(handlers, logger)
	handlers = middleware.SetContext(handlers, logger)
	logger.Infof("Server is starting at %s", conf.Server.GetAddress())
	err = http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
