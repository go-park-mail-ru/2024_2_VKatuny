// Package main starts server and all handlers
package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"
	vacancies_repostory "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	worker_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker/delivery"
	worker_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/worker/repository"
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

	Mux := http.NewServeMux()

	workerRepository := worker_repository.NewRepo()
	workerHandler := worker_delivery.CreateWorkerHandler(workerRepository)
	Mux.Handle("/api/v1/registration/applicant", workerHandler)

	employerRepository := employer_repository.NewRepo()
	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewRepo() // just do it!
	loginHandler := session_delivery.LoginHandler(
		sessionApplicantRepository,
		sessionEmployerRepository,
		workerRepository,
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
	err := http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
