package main

import (
	"log"
	"net/http"
	"places/internal/config"
	"places/internal/repository"
	"places/internal/router"
	"places/internal/storage"
)

func main() {
	logger := config.SetUpLogger("dev")
	logger.Info("Starting the server")

	esStore, err := storage.ConnectWithElasticSearch()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Connected with ElasticSearch")

	data, err := repository.ParsePlacesFromCsv(config.PlacesFilePath, logger)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Parsed places from CSV")

	_, err = esStore.InsertPlaces(data)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Inserted data into ElasticSearch")

	router := router.SetUpRouter(esStore, logger)
	http.ListenAndServe(config.AppPort, router)
}
