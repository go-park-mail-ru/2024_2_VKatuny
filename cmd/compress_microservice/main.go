package main

import (
	"log"
	"net"

	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	//"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	compressdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/delivery"
	compressrepository "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/repository"
	compressusecase "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/compress/usecase"
	compress_api "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"

	"google.golang.org/grpc"
)

func main() {
	//conf, _ := configs.ReadConfig("./configs/conf.yml")
	//logger := logger.NewLogrusLogger()
	lis, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}
	repository := compressrepository.NewCompressRepository("/media/compressed")
	usecase := compressusecase.NewCompressUsecase(repository)
	server := grpc.NewServer()
	compress_api.RegisterCompressServiceServer(server, compressdelivery.NewCompressManager(usecase))

	//logger.Info("starting server at :8091")
	server.Serve(lis)
}
