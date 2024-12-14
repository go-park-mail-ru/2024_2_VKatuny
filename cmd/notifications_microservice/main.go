package main

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	notifications_api "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/generated"
	notificationsdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/delivery"
	notificationsrepository "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/repository"
	notificationsusecase "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/usecase"

	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc"

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

	//go func() {
		lis, err := net.Listen("tcp", conf.NotificationsMicroservice.Server.GetAddress())
		if err != nil {
			log.Fatalln("can't listen port", err)
		}
		repository := notificationsrepository.NewNotificationsRepository(logger, dbConnection)
		usecase := notificationsusecase.NewNotificationsUsecase(repository, logger)
		server := grpc.NewServer()

		notifications_api.RegisterNotificationsServiceServer(server, notificationsdelivery.NewNotificationsManager(usecase, logger))

		logger.Infof("Notifications starting grpc server at %s", conf.NotificationsMicroservice.Server.GetAddress())
		server.Serve(lis)
	// }()

	// http.HandleFunc("/api/v1/notifications/list", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("ws upgrade")
	// 	ws, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	go sendNewMsgNotifications(ws)
	// })
	// logger.Infof("Notifications starting server at 8062", conf.NotificationsMicroservice.Server.GetAddress())
	// http.ListenAndServe(":8062", nil)
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   "e",
		"name":    "n",
		"subject": "s",
	})
	return data
}
