package main

import (
	"fmt"
	"net/http"
	"places/internal/config"
	"places/internal/repository"
	"places/internal/router"
	"places/internal/storage"
)

func main() {
	logger := config.SetUpLogger("dev")
	logger.Info("Starting the server")

	esStore, err := storage.ConnectWithElasticSearch(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while connecting with ES: %s", err))
		return
	}
	logger.Info("Connected with ElasticSearch")

	data, err := repository.ParsePlacesFromCsv(config.PlacesFilePath, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while parsing data from CSV: %s", err))
	}
	logger.Info("Parsed places from CSV")

	_, err = esStore.InsertPlaces(data, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while inserting data: %s", err))
	}
	logger.Info("Inserted data into ElasticSearch")

	router := router.SetUpRouter(esStore, logger)
	http.ListenAndServe(config.AppPort, router)
}
