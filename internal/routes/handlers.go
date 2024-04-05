package routes

import (
	"PlacesApp/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

func PlacesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLOOOOOO")
	switch r.Method {
	case http.MethodGet:
		GetPlaces(w, r)
	case http.MethodPost:
		PostPlace(w, r)
	case http.MethodPut:
		PostPlace(w, r)
	}
}

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.Places)
	fmt.Println("get people")
}

func PostPlace(w http.ResponseWriter, r *http.Request) {
	var place models.Place

	err := json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if err := place.Validate(); err != nil {
		http.Error(w, "Invalid place data", http.StatusBadRequest)
		return
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9201",
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		http.Error(w, "Error setting up Elasticsearch client", http.StatusInternalServerError)
		return
	}

	place = *models.NewPlace(client)
	// if place == nil {
	// 	http.Error(w, "Error creating IndexedPlace", http.StatusInternalServerError)
	// 	return
	// }

	err = place.Index(context.Background(), place)
	if err != nil {
		http.Error(w, "Error indexing place", http.StatusInternalServerError)
		return
	}

	models.Places = append(models.Places, place)

	w.WriteHeader(http.StatusCreated)
}
