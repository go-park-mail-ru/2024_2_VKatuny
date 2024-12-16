package mux

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	cv_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
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
	updateApplicantProfile := middleware.RequireAuthorization(applicantHandlers.UpdateProfile, app, dto.UserTypeApplicant)
	updateApplicantProfile = middleware.CSRFProtection(updateApplicantProfile, app)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/profile", updateApplicantProfile).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/portfolio", applicantHandlers.GetPortfolios).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/cv", applicantHandlers.GetCVs).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/favorite-vacancy", applicantHandlers.GetFavoriteVacancies).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/city", applicantHandlers.GetAllCities).
		Methods(http.MethodGet)

	employerHandlers := employer_delivery.NewEmployerHandlers(app)
	router.HandleFunc("/api/v1/employer/registration", employerHandlers.Registration).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", employerHandlers.GetProfile).
		Methods(http.MethodGet)
	updateEmployerProfile := middleware.RequireAuthorization(employerHandlers.UpdateProfile, app, dto.UserTypeEmployer)
	updateEmployerProfile = middleware.CSRFProtection(updateEmployerProfile, app)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/profile", updateEmployerProfile).
		Methods(http.MethodPut)
	router.HandleFunc("/api/v1/employer/{id:[0-9]+}/vacancies", employerHandlers.GetEmployerVacancies).
		Methods(http.MethodGet)

	cvsHandlers := cv_delivery.NewCVsHandler(app)
	createCV := middleware.RequireAuthorization(cvsHandlers.CreateCV, app, dto.UserTypeApplicant)
	createCV = middleware.CSRFProtection(createCV, app)
	router.HandleFunc("/api/v1/cv", createCV).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", cvsHandlers.GetCV).
		Methods(http.MethodGet)
	updateCV := middleware.RequireAuthorization(cvsHandlers.UpdateCV, app, dto.UserTypeApplicant)
	updateCV = middleware.CSRFProtection(updateCV, app)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", updateCV).
		Methods(http.MethodPut)
	deleteCV := middleware.RequireAuthorization(cvsHandlers.DeleteCV, app, dto.UserTypeApplicant)
	deleteCV = middleware.CSRFProtection(deleteCV, app)
	router.HandleFunc("/api/v1/cv/{id:[0-9]+}", deleteCV).
		Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/cvs", cvsHandlers.SearchCVs).
		Methods(http.MethodGet)
	router.HandleFunc("/api/v1/cv-to-pdf/{id:[0-9]+}", cvsHandlers.CVtoPDF).
		Methods(http.MethodGet)

	vacanciesHandlers := vacancies_delivery.NewVacanciesHandlers(app)
	router.HandleFunc("/api/v1/vacancies", vacanciesHandlers.GetVacancies).
		Methods(http.MethodGet)
	createVacancy := middleware.RequireAuthorization(vacanciesHandlers.CreateVacancy, app, dto.UserTypeEmployer)
	createVacancy = middleware.CSRFProtection(createVacancy, app)
	router.HandleFunc("/api/v1/vacancy", createVacancy).
		Methods(http.MethodPost)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", vacanciesHandlers.GetVacancy).
		Methods(http.MethodGet)
	updateVacancy := middleware.RequireAuthorization(vacanciesHandlers.UpdateVacancy, app, dto.UserTypeEmployer)
	updateVacancy = middleware.CSRFProtection(updateVacancy, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", updateVacancy).
		Methods(http.MethodPut)
	deleteVacancy := middleware.RequireAuthorization(vacanciesHandlers.DeleteVacancy, app, dto.UserTypeEmployer)
	deleteVacancy = middleware.CSRFProtection(deleteVacancy, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}", deleteVacancy).
		Methods(http.MethodDelete)
	AddVacancyIntoFavorite := middleware.RequireAuthorization(vacanciesHandlers.AddVacancyIntoFavorite, app, dto.UserTypeApplicant)
	AddVacancyIntoFavorite = middleware.CSRFProtection(AddVacancyIntoFavorite, app)
	router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/favorite-vacancy", AddVacancyIntoFavorite).
		Methods(http.MethodPost)
	DellVacancyFromFavorite := middleware.RequireAuthorization(vacanciesHandlers.DellVacancyFromFavorite, app, dto.UserTypeApplicant)
	DellVacancyFromFavorite = middleware.CSRFProtection(DellVacancyFromFavorite, app)
		router.HandleFunc("/api/v1/applicant/{id:[0-9]+}/favorite-vacancy", DellVacancyFromFavorite).
			Methods(http.MethodDelete)

	subscribe := middleware.RequireAuthorization(vacanciesHandlers.SubscribeVacancy, app, dto.UserTypeApplicant)
	subscribe = middleware.CSRFProtection(subscribe, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", subscribe).
		Methods(http.MethodPost)
	unsubscribe := middleware.RequireAuthorization(vacanciesHandlers.UnsubscribeVacancy, app, dto.UserTypeApplicant)
	unsubscribe = middleware.CSRFProtection(unsubscribe, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", unsubscribe).
		Methods(http.MethodDelete)
	subscription := middleware.RequireAuthorization(vacanciesHandlers.GetVacancySubscription, app, dto.UserTypeApplicant)
	subscription = middleware.CSRFProtection(subscription, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscription", subscription).
		Methods(http.MethodGet)
	subscribers := middleware.RequireAuthorization(vacanciesHandlers.GetVacancySubscribers, app, dto.UserTypeEmployer)
	subscribers = middleware.CSRFProtection(subscribers, app)
	router.HandleFunc("/api/v1/vacancy/{id:[0-9]+}/subscribers", subscribers).
		Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	return router
}
