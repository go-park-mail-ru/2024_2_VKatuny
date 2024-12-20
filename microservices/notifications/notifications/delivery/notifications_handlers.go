package notificationmicroservice

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	monolit_dto "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	notificationsinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

type NotificationsHandlers struct {
	logger               *logrus.Entry
	notificationsUsecase notificationsinterfaces.INotificationsUsecase
}

func NewNotificationsHandlers(logger *logrus.Logger, notificationsUsecase notificationsinterfaces.INotificationsUsecase) *NotificationsHandlers {
	logger.Debug("NotificationsHandlers created")
	fmt.Println(3)
	return &NotificationsHandlers{
		logger:               logrus.NewEntry(logger),
		notificationsUsecase: notificationsUsecase,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (nd *NotificationsHandlers) GetAlEmployerNotifications(w http.ResponseWriter, r *http.Request) {
	funcName := "NotificationsDelivery.GetAlEmployerNotifications"
	nd.logger.Debugf("%s: got request: ", funcName)

	currentUser, ok := r.Context().Value(monolit_dto.UserContextKey).(*monolit_dto.UserFromSession)
	if !ok {
		w.Write(nil)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		nd.logger.Errorf("%s: %s", funcName, err)
		return
	}
	go func(client *websocket.Conn, employerID uint64) {
		ticker := time.NewTicker(60 * time.Second)
		for {
			fmt.Println("!1")
			w, err := client.NextWriter(websocket.TextMessage)
			if err != nil {
				nd.logger.Errorf("couldn't write: %s", err)
				break
			}
			notificationsList, err := nd.notificationsUsecase.GetAlEmployerNotifications(employerID)
			if err != nil {
				nd.logger.Errorf("could get notifications: %s", err)
				newMessage(w, nil, http.StatusBadRequest)
				continue
			}
			newMessage(w, notificationsList, http.StatusOK)

			<-ticker.C
			fmt.Println("!2")
		}
	}(ws, currentUser.ID)

	go func(client *websocket.Conn, employerID uint64) {
		ticker := time.NewTicker(3 * time.Second)
		for {
			fmt.Println("?1")
			messageType, r, err := client.NextReader()
			if err == nil {
				fmt.Println("?2")
				buffer, err := io.ReadAll(r)
				if err != nil {
					nd.logger.Errorf("could not read message: %s", err)
					continue
				}
				nd.logger.Debugf("messageType: %d", messageType)
				if messageType != 1 {
					nd.logger.Errorf("wrong type read")
					continue
				}
				notificationID, err := strconv.ParseUint(string(buffer[:]), 10, 64) //buffer
				if err != nil {
					nd.logger.Errorf("could not parse notificationID: %s", err)
					continue
				}
				notificationsList, err := nd.notificationsUsecase.GetAlEmployerNotifications(employerID)
				if err != nil {
					nd.logger.Errorf("could get notifications: %s", err)
					continue
				}
				for _, i := range notificationsList {
					if i.ID == notificationID {
						err = nd.notificationsUsecase.MakeEmployerNotificationRead(notificationID)
						if err != nil {
							nd.logger.Errorf("could not make notification read: %s", err)
						}
						continue
					}
				}
				nd.logger.Errorf("not his notification %d", notificationID)
			} else {
				break
			}
			fmt.Println("?3")
			<-ticker.C
		}
	}(ws, currentUser.ID)
}

func newMessage(w io.WriteCloser, notificationsList []*dto.EmployerNotification, status int) {
	data, err := easyjson.Marshal(&dto.JSONResponse{
		HTTPStatus: status,
		Body:       notificationsList,
	})
	if err != nil {
		w.Write(nil)
		return
	}
	w.Write(data)
	w.Close()
}
