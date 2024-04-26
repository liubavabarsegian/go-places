package main

import (
	"log"
	"net/http"
	"places/internal/config"
	"places/internal/db"
	"places/internal/repository"
	"places/internal/routes"
)

func main() {
	esStore, err := db.ConnectWithElasticSearch()
	if err != nil {
		log.Fatal(err)
	}
	data, err := repository.ParsePlacesFromCsv(config.PlacesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Parsed places from CSV", data)

	_, err = esStore.InsertPlaces(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success upload files")

	router := routes.RegisterRoutes(esStore)
	http.ListenAndServe(config.AppPort, router)
}
