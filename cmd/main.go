package main

import (
	"log"
	"places/config"
	"places/db"
	"places/internal/repository"
)

func main() {
	log.Printf("Started app\n")

	es_store, err := db.ConnectWithElasticSearch()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected with Elastic Search")

	data, err := repository.ParsePlacesFromCsv(config.PlacesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Parsed places from CSV")

	_, err = es_store.InsertPlaces(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success upload files")

	places, num, err := es_store.GetPlaces(1000, 100)
	log.Println(places, num, err)

	config.ConfigServer()
}
