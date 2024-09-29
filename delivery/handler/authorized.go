package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/storage"
)

func AuthorizedHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			storage.SetSecureHeaders(w)
			return
		}
		authorized := fmt.Errorf("no user with session")
		session, err := r.Cookie("session_id1")
		var id uint64
		var userType string
		workerBase := BD.HandlersWorker
		employerBase := BD.HandlersEmployer

		if err == nil && session != nil {
			id, authorized = storage.GetWorkerBySession(&workerBase, session)

			userType = "worker"
			if authorized != nil {
				fmt.Println(authorized)
				id, authorized = storage.GetEmployerBySession(&employerBase, session)
				userType = "employer"
				fmt.Println(authorized)
			}
		}
		storage.SetSecureHeaders(w)
		if authorized == nil {
			w.Write([]byte("{statusCode: 200, {id: " + strconv.Itoa(int(id)) + ", usertype: " + userType + "}}"))
		} else {
			w.Write([]byte("{statusCode: 400}"))
		}

	}
	return http.HandlerFunc(fn)
}
