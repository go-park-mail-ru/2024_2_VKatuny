package mux

import (
	"net/http"

	applicant_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/delivery"
	cv_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/delivery"
	employer_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/delivery"
	session_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/delivery"
	vacancies_delivery "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/delivery"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
)

func Init(app *internal.App) *http.ServeMux {
	mux := http.NewServeMux()
	sessionHandlers := session_delivery.NewSessionHandlers(app)
	mux.HandleFunc("/api/v1/login", sessionHandlers.Login)
	mux.HandleFunc("/api/v1/logout", sessionHandlers.Logout)
	mux.HandleFunc("/api/v1/authorized", sessionHandlers.IsAuthorized)

	applicantHandlers := applicant_delivery.NewApplicantProfileHandlers(app)
	mux.HandleFunc("/api/v1/registration/applicant", applicantHandlers.ApplicantRegistration)
	mux.HandleFunc("/api/v1/applicant/profile/", applicantHandlers.ApplicantProfileHandler)
	mux.HandleFunc("/api/v1/applicant/portfolio/", applicantHandlers.GetApplicantPortfoliosHandler)
	mux.HandleFunc("/api/v1/applicant/cv/", applicantHandlers.GetApplicantCVsHandler)

	employerHandlers := employer_delivery.NewEmployerHandlers(app)
	mux.HandleFunc("/api/v1/registration/employer", employerHandlers.EmployerRegistration)
	mux.HandleFunc("/api/v1/employer/profile/", employerHandlers.EmployerProfileHandler)
	mux.HandleFunc("/api/v1/employer/vacancies/", employerHandlers.GetEmployerVacanciesHandler)

	cvsHandlers := cv_delivery.NewCVsHandler(app)
	mux.HandleFunc("/api/v1/cv/", cvsHandlers.CVsRESTHandler)
	mux.HandleFunc("/api/v1/cvs", cvsHandlers.SearchCVs)

	vacanciesHandlers := vacancies_delivery.NewVacanciesHandlers(app)
	mux.HandleFunc("/api/v1/vacancies", vacanciesHandlers.GetVacancies)
	mux.HandleFunc("/api/v1/vacancy/", vacanciesHandlers.VacanciesRESTHandler)
	mux.HandleFunc("/api/v1/vacancy/subscription/", vacanciesHandlers.VacanciesSubscribeRESTHandler)
	mux.HandleFunc("/api/v1/vacancy/subscribers/", vacanciesHandlers.GetVacancySubscribersHandler)

	return mux
}
