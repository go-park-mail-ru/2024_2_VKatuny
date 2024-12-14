package main

import (
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/metrics"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	compressdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/delivery"
	compressrepository "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/repository"
	compressusecase "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/usecase"
	compress_api "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc"
)

func main() {
	conf := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	lis, err := net.Listen("tcp", conf.CompressMicroservice.Server.GetAddress())
	if err != nil {
		log.Fatalln("can't listen port", err)
	}
	repository := compressrepository.NewCompressRepository(conf.CompressMicroservice.CompressedMediaDir, logger)
	usecase := compressusecase.NewCompressUsecase(repository, logger)

	Metrics := metrics.NewMetrics()
	metrics.InitCompressMetrics(Metrics)
	logger.Info("Metrics initialized")
	
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.MetricsInterceptor(Metrics, logger, middleware.CompressMicroservice)),
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(conf.CompressMicroservice.Server.GetMetricsAddress(), mux)
		logger.Infof("Metrics server started at %s", conf.CompressMicroservice.Server.GetMetricsAddress())
	}()

	compress_api.RegisterCompressServiceServer(server, compressdelivery.NewCompressManager(usecase, logger))

	logger.Infof("Compress starting server at %s", conf.CompressMicroservice.Server.GetAddress())
	server.Serve(lis)
}
