package main

import (
	"fmt"
	"log"
	"places/config"
	"places/internal/repository"
)

func main() {
	log.Printf("Started app\n")

	es := config.ConnectWithElasticSearch()
	log.Println("Connected with Elastic Search\n")
	data, err := repository.ParsePlacesFromCsv("config/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed places from data.csv")

	err = repository.InsertPlacesIntoElastic(es, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted places into Elastic\n")

	config.ConfigServer()
}
