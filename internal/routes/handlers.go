package routes

import (
	"PlacesApp/internal/db/models"
	"PlacesApp/internal/repository"
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
	var place2 *repository.Place

	err := json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		fmt.Println(err)
	}

	// Initialize the Elasticsearch client
	// This is a placeholder. You need to replace this with actual client initialization logic.
	// For example:
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Specify the Elasticsearch server address
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
		// Consider sending an error response here
		return
	}

	// Initialize place2 using the NewPlace function
	// Replace 'client' with the actual Elasticsearch client
	place2 = repository.NewPlace(client)

	// Check if place2 is nil before using it
	if place2 == nil {
		fmt.Println("place2 is nil")
		// Consider sending an error response here
		return
	}

	// Now you can safely call the Index method
	err = place2.Index(context.Background(), place)
	if err != nil {
		fmt.Println(err)
		// Consider sending an error response here
		return
	}

	place2.Index(context.Background(), place)

	models.Places = append(models.Places, place)
	fmt.Println("post people")
}
