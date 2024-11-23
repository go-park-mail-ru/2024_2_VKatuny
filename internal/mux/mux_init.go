package mux

import (
	"net/http"

	"github.com/gorilla/mux"

	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	cv_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
)

func Init(app *internal.App) *mux.Router {
	router := mux.NewRouter()
	
	sessionHandlers := session_delivery.NewSessionHandlers(app)
	router.HandleFunc("/api/v1/login", sessionHandlers.Login).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/logout", sessionHandlers.Logout).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/authorized", sessionHandlers.IsAuthorized).
		Methods(http.MethodGet)

	applicantHandlers := applicant_delivery.NewApplicantProfileHandlers(app)
	router.HandleFunc("/api/v1/applicant/registration", applicantHandlers.ApplicantRegistration).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/profile", applicantHandlers.GetProfile).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/profile", applicantHandlers.UpdateProfile).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/portfolio", applicantHandlers.GetPortfolios).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/cv", applicantHandlers.GetCVs).
		Methods(http.MethodGet)

	
	employerHandlers := employer_delivery.NewEmployerHandlers(app)
	router.HandleFunc("/api/v1/employer/registration", employerHandlers.Registration).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", employerHandlers.GetProfile).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", employerHandlers.UpdateProfile).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/vacancies", employerHandlers.GetEmployerVacancies).
		Methods(http.MethodGet)

	cvsHandlers := cv_delivery.NewCVsHandler(app)
	router.HandleFunc("/api/v1/cv", cvsHandlers.CreateCV).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", cvsHandlers.GetCV).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", cvsHandlers.UpdateCV).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", cvsHandlers.DeleteCV).
		Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/cvs", cvsHandlers.SearchCVs).
		Methods(http.MethodGet)

	vacanciesHandlers := vacancies_delivery.NewVacanciesHandlers(app)
	router.HandleFunc("/api/v1/vacancies", vacanciesHandlers.GetVacancies).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/vacancy", vacanciesHandlers.CreateVacancy).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", vacanciesHandlers.GetVacancy).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", vacanciesHandlers.UpdateVacancy).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", vacanciesHandlers.DeleteVacancy).
		Methods(http.MethodDelete)

	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", vacanciesHandlers.SubscribeVacancy).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", vacanciesHandlers.UnsubscribeVacancy).
		Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", vacanciesHandlers.GetVacancySubscription).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscribers", vacanciesHandlers.GetVacancySubscribers).
		Methods(http.MethodGet)

	return router
}
