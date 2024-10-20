// Package main starts server and all handlers
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
	"github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/middleware"
	employer_delivery   "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/delivery"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/clean-arch/internal/pkg/employer/repository"
)

// @title   uArt's API
// @version 1.0

// @contact.name Ifelsik
// @contact.url  https://github.com/Ifelsik

// @host     127.0.0.1:8000
// @BasePath /api/v1
func main() {

	inmemorydb.MakeVacancies()

	inmemorydb.MakeUsers()
	Mux := http.NewServeMux()

	workerHandler := worker_delivery.CreateWorkerHandler(&inmemorydb.HandlersWorker)
	workerHandler = middleware.SetSecurityAndOptionsHeaders(workerHandler)
	Mux.Handle("/api/v1/registration/worker", workerHandler)

	employerRepository := employer_repository.NewRepo()
	employerHandler := employer_delivery.CreateEmployerHandler(employerRepository)
	employerHandler =  middleware.SetSecurityAndOptionsHeaders(employerHandler)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	// change handler's destination
	// loginHandler := handler.LoginHandler()
	// Mux.Handle("/api/v1/login", loginHandler)

	// change handler's destination
	// logoutHandler := handler.LogoutHandler()
	// Mux.Handle("/api/v1/logout", logoutHandler)

	// change handler's destination
	// authorizedHandler := handler.AuthorizedHandler()
	// Mux.Handle("/api/v1/authorized", authorizedHandler)

	// change handler's destination
	vacanciesListHandler := vacanciesDelivery.VacanciesHandler() //(&db.Vacancies)
	vacanciesListHandler =  middleware.SetSecurityAndOptionsHeaders(vacanciesListHandler)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

	log.Print("Listening...")
	http.ListenAndServe(inmemorydb.BACKENDIP, Mux)
	fmt.Print("started")
}
