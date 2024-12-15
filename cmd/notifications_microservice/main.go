package main

import (
	"log"
	"net"

	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	grpc_auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	notifications_api "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"
	notificationsdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/delivery"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/mux"
	notificationsrepository "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/repository"
	notificationsusecase "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendNewMsgNotifications(client *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			ticker.Stop()
			break
		}

		msg := newMessage()
		w.Write(msg)
		w.Close()

		<-ticker.C
	}
}
func main() {
	conf := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()
	dbConnection, err := utils.GetDBConnection(conf.DataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()
	repository := notificationsrepository.NewNotificationsRepository(logger, dbConnection)
	usecase := notificationsusecase.NewNotificationsUsecase(repository, logger)
	go func() {
		lis, err := net.Listen("tcp", conf.NotificationsMicroservice.GRPCserver.GetAddress())
		if err != nil {
			log.Fatalln("can't listen port", err)
		}

		server := grpc.NewServer()

		notifications_api.RegisterNotificationsServiceServer(server, notificationsdelivery.NewNotificationsManager(usecase, logger))

		logger.Infof("Notifications starting grpc server at %s", conf.NotificationsMicroservice.GRPCserver.GetAddress())
		server.Serve(lis)
	}()

	connAuthGRPC, err := grpc.NewClient(
		conf.AuthMicroservice.Server.GetAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer connAuthGRPC.Close()
	microservices := &internal.Microservices{
		Auth: grpc_auth.NewAuthorizationClient(connAuthGRPC),
	}
	logger.Infof("gRPC client started at %s", conf.AuthMicroservice.Server.GetAddress())
	app := &internal.App{
		Logger:        logger,
		Microservices: microservices,
	}

	Mux := mux.Init(app, logger, usecase)


	handlers := middleware.SetSecurityAndOptionsHeaders(Mux, conf.Server.Front)
	// http.HandleFunc("/api/v1/notifications/list", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("ws upgrade")
	// 	ws, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	go sendNewMsgNotifications(ws)
	// })
	logger.Infof("Notifications starting server at %s", conf.NotificationsMicroservice.Server.GetAddress())
	err = http.ListenAndServe(conf.NotificationsMicroservice.Server.GetAddress(), handlers)
	if err != nil {
		log.Fatal(err)
	}
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   "e",
		"name":    "n",
		"subject": "s",
	})
	return data
}
