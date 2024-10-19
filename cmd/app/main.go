// Package main starts server and all handlers
package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/article/delivery/handler"
	"github.com/go-park-mail-ru/2024_2_VKatuny/inmemorydb"
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

	workerHandler := handler.CreateWorkerHandler(&inmemorydb.HandlersWorker)
	Mux.Handle("/api/v1/registration/worker", workerHandler)

	employerHandler := handler.CreateEmployerHandler(&inmemorydb.HandlersEmployer)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	loginHandler := handler.LoginHandler()
	Mux.Handle("/api/v1/login", loginHandler)

	logoutHandler := handler.LogoutHandler()
	Mux.Handle("/api/v1/logout", logoutHandler)

	authorizedHandler := handler.AuthorizedHandler()
	Mux.Handle("/api/v1/authorized", authorizedHandler)

	vacanciesListHandler := handler.VacanciesHandler() //(&db.Vacancies)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

	log.Print("Listening...")
	http.ListenAndServe(inmemorydb.BACKENDIP, Mux)
}
