package config

import (
	"PlacesApp/internal/routes"
	"context"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
)

type ClientKey struct{}

var CTX context.Context

func ConnectWithElasticSearch(ctx context.Context) context.Context {
	es_config := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9201",
		},
	}
	newClient, err := elasticsearch.NewClient(es_config)
	if err != nil {
		panic(err)
	}

	return context.WithValue(ctx, ClientKey{}, newClient)
}

func ConfigServer(ctx context.Context) {
	CTX = ctx
	r := mux.NewRouter()
	routes.RegisterPlacesAppRoutes(r)
	log.Fatal(http.ListenAndServe("localhost:9200", r))
}
