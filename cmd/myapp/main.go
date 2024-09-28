package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/BD"
	"github.com/go-park-mail-ru/2024_2_VKatuny/delivery/handler"
)

func main() {
	Mux := http.NewServeMux()

	workerHandler := handler.CreateWorkerHandler(&BD.HandlersWorker)
	Mux.Handle("/api/registration/worker", workerHandler)
	employerHandler := handler.CreateEmployerHandler(&BD.HandlersEmployer)
	Mux.Handle("/api/registration/employer", employerHandler)

	log.Print("Listening...")
	http.ListenAndServe("0.0.0.0:8080", Mux)

	// fmt.Println(EqualHashedPasswords(HashPassword("pass"), "pass"))
	// http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte("{}"))
	// })

	// http.HandleFunc("/api/registration/worker", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")

	// 	log.Println(r.URL.Path)
	// 	//fmt.Println("11", r.Method)
	// 	if r.Method == http.MethodPost {
	// 		fmt.Println("2")
	// 		HandlersWorker.HandleCreateWorker(w, r)
	// 		return
	// 	}

	// })

	// http.HandleFunc("/api/registration/employer", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")

	// 	log.Println(r.URL.Path)
	// 	fmt.Println("11", r.Method)
	// 	if r.Method == http.MethodPost {
	// 		fmt.Println("3")
	// 		HandlersEmployer.HandleCreateEmployer(w, r)
	// 		return
	// 	}

	// })
	// http.ListenAndServe("0.0.0.0:8080", nil)
}
