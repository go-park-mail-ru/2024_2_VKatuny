package main

import (
	"net"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth"
	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	conf := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	pgSQLConn, err := utils.GetDBConnection(conf.DataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Successfully connected to postgres")
	defer pgSQLConn.Close()

	redisConn, err := redis.Dial("tcp", conf.AuthMicroservice.Database.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer redisConn.Close()

	if _, err := redisConn.Do("AUTH", conf.AuthMicroservice.Database.User, conf.AuthMicroservice.Database.Password); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Successfully connected to redis")

	lister, err := net.Listen("tcp", conf.AuthMicroservice.Server.GetAddress())
	if err != nil {
		logger.Fatalf("failed to listen port: %s", err)
	}

	server := grpc.NewServer()

	grpc_auth.RegisterAuthorizationServer(server, auth.NewAuthorization(pgSQLConn, redisConn, logger))

	logger.Infof("Starting gRPC server on %s", conf.AuthMicroservice.Server.GetAddress())
	server.Serve(lister)
}
