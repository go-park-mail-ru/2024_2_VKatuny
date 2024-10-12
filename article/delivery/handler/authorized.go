package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/article/repository"
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
		var typeOfUser string
		workerBase := BD.HandlersWorker
		employerBase := BD.HandlersEmployer

		if err == nil && session != nil {
			id, authorizationErr = repository.GetWorkerBySession(&workerBase, session)

			typeOfUser = BD.WORKER
			if authorizationErr != nil {
				log.Println(authorizationErr)
				id, authorizationErr = repository.GetEmployerBySession(&employerBase, session)
				typeOfUser = BD.EMPLOYER
				log.Println(authorizationErr)
			}
		}

		if authorizationErr == nil {
			repository.UniversalMarshal(w, http.StatusOK, BD.ReturnUserFields{200, BD.AuthorizedUserFields{id, typeOfUser}})
		} else {
			repository.UniversalMarshal(w, http.StatusUnauthorized, nil)
		}
	}))
}
