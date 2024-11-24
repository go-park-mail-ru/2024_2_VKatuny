// Package main starts server and all handlers
package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"

	applicant_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	cvRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	file_loading_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/repository"
	file_loading_usecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/usecase"
	portfolioRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"

	_ "github.com/jackc/pgx/v5/stdlib"

	vacancies_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
)

var (
	CompressManager compressmicroservice.CompressServiceClient
)

// @title   μArt's API
// @version 1.0

// @contact.name Ifelsik
// @contact.url  https://github.com/Ifelsik
// @contact.name Olgmuzalev13
// @contact.url  https://github.com/Olgmuzalev13

// @host     127.0.0.1:8080
// @BasePath /api/v1
func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	dbConnection, err := utils.GetDBConnection(conf.DataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()
	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewSessionStorage(dbConnection)
	repositories := &internal.Repositories{
		ApplicantRepository:        applicant_repository.NewApplicantStorage(dbConnection),
		PortfolioRepository:        portfolioRepository.NewPortfolioStorage(dbConnection),
		CVRepository:               cvRepository.NewCVStorage(dbConnection),
		VacanciesRepository:        vacancies_repository.NewVacanciesStorage(dbConnection),
		EmployerRepository:         employer_repository.NewEmployerStorage(dbConnection),
		SessionApplicantRepository: sessionApplicantRepository,
		SessionEmployerRepository:  sessionEmployerRepository,
		FileLoadingRepository:      file_loading_repository.NewFileLoadingStorage(conf.Server.Front),
	}
	usecases := &internal.Usecases{

		ApplicantUsecase:   applicantUsecase.NewApplicantUsecase(logger, repositories, CompressManager),
		FileLoadingUsecase: file_loading_usecase.NewFileLoadingUsecase(logger, repositories),
	}
	app := &internal.App{
		Logger:       logger,
		Repositories: repositories,
		Usecases:     usecases,
	}

	Mux := mux.Init(app)

	// Wrapped multiplexer
	// Mux implements http.Handler interface so it's possible to wrap
	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.MediaDir)
	handlers = middleware.AccessLogger(handlers, logger)
	handlers = middleware.SetLogger(handlers, logger)
	handlers = middleware.Panic(handlers, logger)
	logger.Infof("Server is starting at %s", conf.Server.GetAddress())
	err = http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
