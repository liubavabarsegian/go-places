package config

import (
	"PlacesApp/internal/routes"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
)

func ConnectWithElasticSearch() *elasticsearch.Client {
	es_config := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
	}
	newClient, err := elasticsearch.NewClient(es_config)
	if err != nil {
		panic(err)
	}

	return newClient
}

func ConfigServer() {
	r := mux.NewRouter()
	routes.RegisterPlacesAppRoutes(r)
	http.Handle("/", r)
	fmt.Println("config")
	http.ListenAndServe("127.0.0.1:8888", r)
}
