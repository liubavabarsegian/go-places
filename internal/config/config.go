package config

import (
	"net/http"
	"places/internal/routes"
)

const (
	PlacesFilePath = "/app/internal/config/data.csv"
	ElasticAddress = "http://elasticsearch:9200"
	IndexName      = "places"
	Schema         = "/app/internal/config/schema.json"
	AppPort        = ":8888"
)

func ConfigServer() {
	router := routes.RegisterRoutes()
	http.ListenAndServe(AppPort, router)
}
