package notificationmicroservice

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	monolit_dto "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	notificationsinterfaces "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/dto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type NotificationsHandlers struct {
	logger               *logrus.Entry
	notificationsUsecase notificationsinterfaces.INotificationsUsecase
}

func NewNotificationsHandlers(logger *logrus.Logger, notificationsUsecase notificationsinterfaces.INotificationsUsecase) *NotificationsHandlers {
	logger.Debug("NotificationsHandlers created")
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
	ticker := time.NewTicker(3 * time.Second)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		nd.logger.Errorf("%s: %s", funcName, err)
	}
	go func(client *websocket.Conn, employerID uint64) {
		for {
			w, err := client.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			defer w.Close()
			if false {
				newMessage(w, nil)
				return
			}
			notificationsList, err := nd.notificationsUsecase.GetAlEmployerNotifications(employerID)
			if err != nil {
				newMessage(w, nil)
				return
			}
			newMessage(w, notificationsList)
			<-ticker.C
		}
	}(ws, currentUser.ID)
}

func newMessage(w io.WriteCloser, notificationsList []*dto.EmployerNotification) {
	data, err := json.Marshal(notificationsList)
	if err != nil {
		w.Write(nil)
		return
	}
	w.Write(data)
}

// func (nd *NotificationsHandlers) MakeEmployerNotificationRead(w http.ResponseWriter, r *http.Request) {
// 	funcName := "NotificationsDelivery.MakeEmployerNotificationRead"
// 	nd.logger.Debugf("%s: got request: ", funcName)
// 	vars := mux.Vars(r)
// 	slug := vars["id"]
// 	notificationID, err := strconv.ParseUint(slug, 10, 64)
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		nd.logger.Errorf("%s: %s", funcName, err)
// 	}
// 	go func(client *websocket.Conn, notificationID uint64) {
// 		if false {
// 			nd.logger.Errorf("bad input: %s", false)
// 			middleware.UniversalMarshal(w, http.StatusBadRequest, dto.JSONResponse{
// 				HTTPStatus: http.StatusBadRequest,
// 				Error:      dto.MsgInvalidJSON,
// 			})
// 			return
// 		}
// 		err := nd.notificationsUsecase.MakeEmployerNotificationRead(notificationID)
// 		if err != nil {
// 			middleware.UniversalMarshal(w, http.StatusInternalServerError, dto.JSONResponse{
// 				HTTPStatus: http.StatusInternalServerError,
// 				Error:      err.Error(),
// 			})
// 			return
// 		}

// 		middleware.UniversalMarshal(w, http.StatusOK, dto.JSONResponse{
// 			HTTPStatus: http.StatusOK,
// 			Body:       nil,
// 		})
// 	}(ws, notificationID)
// }
