package routes

import (
	"PlacesApp/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

type ClientKey struct{}

func PlacesHandler(w http.ResponseWriter, r *http.Request) {
	// client, ok := r.Context().Value(ClientKey{}).(*elasticsearch.Client)
	// if !ok {
	// 	http.Error(w, "Elasticsearch client not found in context", http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Println("places handler", client.API)
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
	fmt.Println("get places")
}

func PostPlace(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post places")
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

	client, ok := r.Context().Value(ClientKey{}).(*elasticsearch.Client)
	if !ok {
		http.Error(w, "Elasticsearch client not found in context 2", http.StatusInternalServerError)
		return
	}
	// if err != nil {
	// 	http.Error(w, "Error setting up Elasticsearch client", http.StatusInternalServerError)
	// 	return
	// }

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
