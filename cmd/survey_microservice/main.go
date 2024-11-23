package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/configs"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/logger"
	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/go-park-mail-ru/2024_2_VKatuny/microservices/survey/survey/delivery"

	_ "github.com/jackc/pgx/v5/stdlib"

)

func main() {
	conf, _ := configs.ReadConfig("./configs/conf.yml")
	logger := logger.NewLogrusLogger()

	dbConnection, err := utils.GetDBConnection(conf.SurveyDataBase.GetDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbConnection.Close()

	router := mux.NewRouter()

	handlers := delivery.NewSurveyHandlers(logger)

	router.HandleFunc("/api/v1/survey/statistics", handlers.GetStatistics).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/survey", handlers.GetSurveyForm).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/survey", handlers.AddSurveyAnswer).Methods(http.MethodPost)

	logger.Infof("Server is starting at %s", conf.Server.MicroserviceSurveyURI)
	err = http.ListenAndServe(conf.Server.MicroserviceSurveyURI, router)
	if err != nil {
		logger.Fatal(err)
	}
}
