package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/delivery/handler"
)

func main() {
	Mux := http.NewServeMux()

	workerHandler := handler.CreateWorkerHandler(&BD.HandlersWorker)
	Mux.Handle("/api/registration/worker/", workerHandler)

	employerHandler := handler.CreateEmployerHandler(&BD.HandlersEmployer)
	Mux.Handle("/api/registration/employer/", employerHandler)

	loginHandler := handler.LoginHandler()
	Mux.Handle("/api/login", loginHandler)

	logoutHandler := handler.LogoutHandler()
	Mux.Handle("/api/logout", logoutHandler)

	authorizedHandler := handler.AuthorizedHandler()
	Mux.Handle("/api/authorized", authorizedHandler)

	log.Print("Listening...")
	http.ListenAndServe("0.0.0.0:8080", Mux)
}
