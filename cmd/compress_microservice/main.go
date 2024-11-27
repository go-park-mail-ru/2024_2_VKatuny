package main

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	compressdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/delivery"
	compressrepository "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/repository"
	compressusecase "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/usecase"
	compress_api "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"

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
	server := grpc.NewServer()
	compress_api.RegisterCompressServiceServer(server, compressdelivery.NewCompressManager(usecase, logger))

	logger.Infof("Compress starting server at %s", conf.CompressMicroservice.Server.GetAddress())
	server.Serve(lis)
}
