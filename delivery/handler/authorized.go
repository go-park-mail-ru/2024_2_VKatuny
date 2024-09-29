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
		defer r.Body.Close()
		isoption := storage.Isoption(w, r)
		if isoption {
			return
		}
		storage.SetSecureHeaders(w)
		w.Header().Set("Content-Type", "application/json")
		authorized := fmt.Errorf("no user with session")
		session, err := r.Cookie("session_id1")
		var id uint64
		var userType string
		workerBase := BD.HandlersWorker
		employerBase := BD.HandlersEmployer

		fmt.Println(BD.HandlersEmployer)

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

		if authorized == nil {
			w.Write([]byte("{statusCode: 200, {id: " + strconv.Itoa(int(id)) + ", usertype: " + userType + "}}"))
		} else {
			w.WriteHeader(401)
		}
	}
	return http.HandlerFunc(fn)
}
