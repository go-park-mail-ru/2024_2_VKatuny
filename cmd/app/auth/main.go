package auth

import (
	"net"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth"
	grpc_server "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"google.golang.org/grpc"
)

func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()
	
	dbConnection, err := utils.GetDBConnection(conf.DataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	lister, err := net.Listen("tcp", ":8091")
	if err != nil {
		logger.Fatalf("failed to listen port: %s", err)
	}

	server := grpc.NewServer()

	grpc_server.RegisterAuthorizationServer(server, auth.NewAuthorization(dbConnection, logger))

	logger.Info("Starting gRPC server on port 8091")
	server.Serve(lister)
}
