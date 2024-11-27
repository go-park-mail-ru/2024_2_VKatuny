// Package main starts server and all handlers
package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	applicant_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/repository"
	applicantUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/applicant/usecase"
	cvRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/repository"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"
	employer_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/repository"
	file_loading_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/repository"
	file_loading_usecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/file_loading/usecase"
	portfolioRepository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/repository"
	session_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"
	session_usecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/usecase"
	vacanciesUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"

	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"

	_ "github.com/jackc/pgx/v5/stdlib"

	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	vacancies_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"

	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"
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
	conf := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	dbConnection, err := utils.GetDBConnection(conf.DataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	connAuthGRPC, err := grpc.NewClient(
		conf.Server.GetAuthServiceLocation(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer connAuthGRPC.Close()
	logger.Infof("gRPC client started at %s", conf.AuthMicroservice.Server.GetAddress())

	connCompressGRPC, err := grpc.NewClient(
		conf.CompressMicroservice.Server.GetAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer connAuthGRPC.Close()
	logger.Infof("Compress gRPC client started at %s", conf.CompressMicroservice.Server.GetAddress())

	sessionApplicantRepository, sessionEmployerRepository := session_repository.NewSessionStorage(dbConnection)
	repositories := &internal.Repositories{
		ApplicantRepository:        applicant_repository.NewApplicantStorage(dbConnection),
		PortfolioRepository:        portfolioRepository.NewPortfolioStorage(dbConnection),
		CVRepository:               cvRepository.NewCVStorage(dbConnection),
		VacanciesRepository:        vacancies_repository.NewVacanciesStorage(dbConnection),
		EmployerRepository:         employer_repository.NewEmployerStorage(dbConnection),
		SessionApplicantRepository: sessionApplicantRepository,
		SessionEmployerRepository:  sessionEmployerRepository,
		FileLoadingRepository:      file_loading_repository.NewFileLoadingStorage(conf.Server.MediaDir),
	}
	microservices := &internal.Microservices{
		Auth:     grpc_auth.NewAuthorizationClient(connAuthGRPC),
		Compress: compressmicroservice.NewCompressServiceClient(connCompressGRPC),
	}
	usecases := &internal.Usecases{
		ApplicantUsecase:   applicantUsecase.NewApplicantUsecase(logger, repositories),
		PortfolioUsecase:   portfolioUsecase.NewPortfolioUsecase(logger, repositories),
		CVUsecase:          cvUsecase.NewCVsUsecase(logger, repositories),
		VacanciesUsecase:   vacanciesUsecase.NewVacanciesUsecase(logger, repositories),
		EmployerUsecase:    employerUsecase.NewEmployerUsecase(logger, repositories),
		SessionUsecase:     session_usecase.NewSessionUsecase(logger, repositories),
		FileLoadingUsecase: file_loading_usecase.NewFileLoadingUsecase(logger, repositories, microservices, conf),
	}

	app := &internal.App{
		Logger:        logger,
		Repositories:  repositories,
		Usecases:      usecases,
		Microservices: microservices,
	}

	Mux := mux.Init(app)

	// Wrapped multiplexer
	// Mux implements http.Handler interface so it's possible to wrap
	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.Front)
	handlers = middleware.AccessLogger(handlers, logger)
	handlers = middleware.SetLogger(handlers, logger)
	handlers = middleware.Panic(handlers, logger)
	logger.Infof("Server is starting at %s", conf.Server.GetAddress())
	err = http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
