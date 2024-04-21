package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"places/internal/entities"
)

type ClientKey struct{}

func PlacesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetPlaces(w, r)
	case http.MethodPost:
		// PostPlace(w, r)
	}
}

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entities.Places)
	fmt.Println("get places")
}

// func PostPlace(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("post places")
// 	var place entities.Place

// 	err := json.NewDecoder(r.Body).Decode(&place)
// 	if err != nil {
// 		http.Error(w, "Error decoding request body", http.StatusBadRequest)
// 		return
// 	}

// 	if err := place.Validate(); err != nil {
// 		http.Error(w, "Invalid place data", http.StatusBadRequest)
// 		return
// 	}

// 	client, ok := r.Context().Value(ClientKey{}).(*elasticsearch.Client)
// 	if !ok {
// 		http.Error(w, "Elasticsearch client not found in context 2", http.StatusInternalServerError)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, "Error setting up Elasticsearch client", http.StatusInternalServerError)
// 		return
// 	}

// 	place = *entities.NewPlace(client)
// 	// if place == nil {
// 	// 	http.Error(w, "Error creating IndexedPlace", http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	err = place.Index(context.Background(), place)
// 	if err != nil {
// 		http.Error(w, "Error indexing place", http.StatusInternalServerError)
// 		return
// 	}

// 	entities.Places = append(entities.Places, place)

// 	w.WriteHeader(http.StatusCreated)
// }
