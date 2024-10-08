package main

import (
	"log"
	"net/http"
	
	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/delivery/handler"
)

// @title   uArt's API
// @version 1.0

// @contact.name Ifelsik
// @contact.url  https://github.com/Ifelsik

// @host     127.0.0.1:8000
// @BasePath /api/v1
func main() {
	BD.MakeVacancies()

	BD.MakeUsers()
	Mux := http.NewServeMux()

	workerHandler := handler.CreateWorkerHandler(&BD.HandlersWorker)
	Mux.Handle("/api/v1/registration/worker", workerHandler)

	employerHandler := handler.CreateEmployerHandler(&BD.HandlersEmployer)
	Mux.Handle("/api/v1/registration/employer", employerHandler)

	loginHandler := handler.LoginHandler()
	Mux.Handle("/api/v1/login", loginHandler)

	logoutHandler := handler.LogoutHandler()
	Mux.Handle("/api/v1/logout", logoutHandler)

	authorizedHandler := handler.AuthorizedHandler()
	Mux.Handle("/api/v1/authorized", authorizedHandler)

	vacanciesListHandler := handler.VacanciesHandler(&BD.Vacancies)
	Mux.Handle("/api/v1/vacancies", vacanciesListHandler)

	log.Print("Listening...")
	http.ListenAndServe(BD.BACKENDIP, Mux)
}
