package config

import (
	"net/http"
	"places/internal/routes"
)

const (
	PlacesFilePath = "/app/config/data.csv"
	ElasticAddress = "http://elasticsearch:9200"
	IndexName      = "places"
	Schema         = "/app/config/schema.json"
)

func ConfigServer() {
	router := routes.RegisterRoutes()
	http.ListenAndServe(":8888", router)
}
