package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/survey/survey/delivery"
)

func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	router := mux.NewRouter()

	handlers := delivery.NewSurveyHandlers(logger)

	router.HandleFunc("/api/v1/survey/statistics", handlers.GetStatistics).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/survey", handlers.GetSurveyForm).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/survey", handlers.AddSurveyAnswer).Methods(http.MethodPost)

	http.ListenAndServe(conf.Server.MicroserviceSurveyURI, router)
}
