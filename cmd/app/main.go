// Package main starts server and all handlers
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	portfolioUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/portfolio/usecase"

	vacanciesUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/usecase"

	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	notificationsmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/mux"

	_ "github.com/jackc/pgx/v5/stdlib"

	cvUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/cvs/usecase"
	vacancies_repository "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/vacancies/repository"
	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"

	employerUsecase "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/employer/usecase"
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
	defer connCompressGRPC.Close()
	logger.Infof("Compress gRPC client started at %s", conf.CompressMicroservice.Server.GetAddress())
	connNotificationsGRPC, err := grpc.NewClient(
		conf.NotificationsMicroservice.GRPCserver.GetAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer connNotificationsGRPC.Close()
	logger.Infof("Notifications gRPC client started at %s", conf.NotificationsMicroservice.GRPCserver.GetAddress())

	Metrics := metrics.NewMetrics()
	metrics.InitMetrics(Metrics)

	repositories := &internal.Repositories{
		ApplicantRepository:        applicant_repository.NewApplicantStorage(dbConnection, logger),
		PortfolioRepository:        portfolioRepository.NewPortfolioStorage(dbConnection, logger),
		CVRepository:               cvRepository.NewCVStorage(dbConnection, logger),
		VacanciesRepository:        vacancies_repository.NewVacanciesStorage(dbConnection, logger),
		EmployerRepository:         employer_repository.NewEmployerStorage(dbConnection, logger),
		FileLoadingRepository:      file_loading_repository.NewFileLoadingStorage(logger, conf.Server.MediaDir, conf.Server.CVinPDFDir, conf.Server.TamplateDir),
	}
	microservices := &internal.Microservices{
		Auth:     grpc_auth.NewAuthorizationClient(connAuthGRPC),
		Compress: compressmicroservice.NewCompressServiceClient(connCompressGRPC),
		Notifications: notificationsmicroservice.NewNotificationsServiceClient(connNotificationsGRPC),
	}
	usecases := &internal.Usecases{
		ApplicantUsecase:   applicantUsecase.NewApplicantUsecase(logger, repositories),
		PortfolioUsecase:   portfolioUsecase.NewPortfolioUsecase(logger, repositories),
		CVUsecase:          cvUsecase.NewCVsUsecase(logger, repositories),
		VacanciesUsecase:   vacanciesUsecase.NewVacanciesUsecase(logger, repositories),
		EmployerUsecase:    employerUsecase.NewEmployerUsecase(logger, repositories),
		FileLoadingUsecase: file_loading_usecase.NewFileLoadingUsecase(logger, repositories, microservices, conf),
	}

	app := &internal.App{
		Logger:        logger,
		Repositories:  repositories,
		CSRFSecret:    conf.Server.CSRFSecret,
		Usecases:      usecases,
		Microservices: microservices,
		Metrics:       Metrics,
	}

	Mux := mux.Init(app)

	Mux.Handle("/metrics", promhttp.Handler())

	// Wrapped multiplexer
	// Mux implements http.Handler interface so it's possible to wrap
	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.Front)
	handlers = middleware.AccessLogger(handlers, logger, app.Metrics)
	handlers = middleware.SetLogger(handlers, logger)
	handlers = middleware.Panic(handlers, logger)

	logger.Debugf("Starting compress demon")
	_, err = microservices.Compress.StartScanCompressDemon(
		context.Background(),
		&compressmicroservice.Nothing{},
	)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Server is starting at %s", conf.Server.GetAddress())
	err = http.ListenAndServe(conf.Server.GetAddress(), handlers)
	if err != nil {
		logger.Fatal(err)
	}
}
