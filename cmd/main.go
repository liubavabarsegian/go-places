package main

import (
	"PlacesApp/config"
	"PlacesApp/internal/domain"
	"fmt"
	"log"
)

func main() {
	es := config.ConnectWithElasticSearch()

	data, err := domain.ParseDataFromCsv("config/data.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("parsed")
	// res, err := http.Get("http://localhost:9200")

	// if err != nil {
	// 	log.Println("couldnt send get request")
	// }

	// if res != nil {
	// 	defer res.Body.Close()
	// 	fmt.Println("Response status:", res.Status)
	// }

	// res, err = http.Get("http://elasticsearch:9200")
	// if err != nil {
	// 	log.Println("couldnt send elastic get request")
	// }
	// if res != nil {
	// 	defer res.Body.Close()
	// 	fmt.Println("Response status:", res.Status)
	// }

	// log.Println(res)
	err = domain.InsertDataToElastic(es, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted")

	config.ConfigServer()
}
