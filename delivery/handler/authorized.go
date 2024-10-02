package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

// Authorized godoc
// @Summary     Checks user's authorization 
// @Description Gets cookie from user and checks authentication
// @Tags        AuthStatus
// @Param       session_id header string true "Session ID (Cookie)"
// @Success     200
// @Failure     401
// @Router      /authorized [post]
func AuthorizedHandler() http.Handler {
	return HttpHeadersWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		
		authorizationErr := fmt.Errorf("no user with session")
		session, err := r.Cookie("session_id1")
		var id uint64
		var userType string
		workerBase := BD.HandlersWorker
		employerBase := BD.HandlersEmployer

		if err == nil && session != nil {
			id, authorizationErr = storage.GetWorkerBySession(&workerBase, session)

			userType = "worker"
			if authorizationErr != nil {
				log.Println(authorizationErr)
				id, authorizationErr = storage.GetEmployerBySession(&employerBase, session)
				userType = "employer"
				log.Println(authorizationErr)
			}
		}

		if authorizationErr == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"statusCode": 200, "user": {"id": ` + strconv.Itoa(int(id)) + `, "usertype": "` + userType + `"}}`))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
}
