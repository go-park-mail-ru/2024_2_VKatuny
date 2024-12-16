package mux

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/commonerrors"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	notificationsmicroserviceinterface "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications"
	notificationsdelivery "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/notifications/notifications/delivery"
)

func Init(app *internal.App, logger *logrus.Logger, notificationsUsecase notificationsmicroserviceinterface.INotificationsUsecase) *mux.Router {
	fmt.Println(2)
	router := mux.NewRouter()
	notificationsHandlers := notificationsdelivery.NewNotificationsHandlers(logger, notificationsUsecase)
	notificationList := middleware.RequireAuthorization(notificationsHandlers.GetAlEmployerNotifications, app, dto.UserTypeEmployer)

	router.HandleFunc("/api/v1/notifications", notificationList)
	return router
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusNotFound)
	response := &dto.JSONResponse{
		HTTPStatus: http.StatusNotFound,
		Error:      commonerrors.ErrFrontServiceNotFound.Error(),
	}
	JSONResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(JSONResponse)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusMethodNotAllowed)
	response := &dto.JSONResponse{
		HTTPStatus: http.StatusMethodNotAllowed,
		Error:      commonerrors.ErrFrontMethodNotAllowed.Error(),
	}
	JSONResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(JSONResponse)
}
