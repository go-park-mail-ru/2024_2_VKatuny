// // Package main starts server and all handlers
// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-park-mail-ru/2024_2_VKatuny/article/delivery/handler"
// 	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
// )

// // @title   uArt's API
// // @version 1.0

// // @contact.name Ifelsik
// // @contact.url  https://github.com/Ifelsik

// // @host     127.0.0.1:8000
// // @BasePath /api/v1
// func main() {

// 	inmemorydb.MakeVacancies()

// 	inmemorydb.MakeUsers()
// 	Mux := http.NewServeMux()

// 	workerHandler := handler.CreateWorkerHandler(&inmemorydb.HandlersWorker)
// 	Mux.Handle("/api/v1/registration/worker", workerHandler)

// 	employerHandler := handler.CreateEmployerHandler(&inmemorydb.HandlersEmployer)
// 	Mux.Handle("/api/v1/registration/employer", employerHandler)

// 	loginHandler := handler.LoginHandler()
// 	Mux.Handle("/api/v1/login", loginHandler)

// 	logoutHandler := handler.LogoutHandler()
// 	Mux.Handle("/api/v1/logout", logoutHandler)

// 	authorizedHandler := handler.AuthorizedHandler()
// 	Mux.Handle("/api/v1/authorized", authorizedHandler)

// 	vacanciesListHandler := handler.VacanciesHandler() //(&db.Vacancies)
// 	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

// 	log.Print("Listening...")
// 	http.ListenAndServe(inmemorydb.BACKENDIP, Mux)
// 	fmt.Print("started")
// }

// Package main starts server and all handlers
package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/inmemorydb"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/repository"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/session/delivery"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/vacancies/delivery"
	vacancies_repostory "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/vacancies/repository"
	worker_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker/delivery"
	worker_repository "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/worker/repository"
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
	// fill db with data using vacancies/repostory
	// inmemorydb.MakeVacancies()

	inmemorydb.MakeUsers()
	Mux := http.NewServeMux()

	workerRepository := worker_repository.NewRepo()
	workerHandler := worker_delivery.CreateWorkerHandler(workerRepository)
	Mux.Handle("/api/v1/registration/applicant", workerHandler)

	employerRepository := employer_repository.NewRepo()
	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	sessionRepository := session_delivery.NewRepo() // just do it!
	loginHandler := session_delivery.LoginHandler(sessionRepository, conf.Server.GetAddress())
	Mux.Handle("/api/v1/login", loginHandler)

	logoutHandler := session_delivery.LogoutHandler(sessionRepository)
	Mux.Handle("/api/v1/logout", logoutHandler)

	authorizedHandler := session_delivery.AuthorizedHandler(sessionRepository)
	Mux.Handle("/api/v1/authorized", authorizedHandler)

	vacanciesRepository := vacancies_repostory.NewRepo()
	vacanciesListHandler := vacancies_delivery.VacanciesHandler(vacanciesRepository) //(&db.Vacancies)
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
