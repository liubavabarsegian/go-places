package routes

import (
	"PlacesApp/internal/db/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func PlacesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLOOOOOO")
	switch r.Method {
	case http.MethodGet:
		GetPlaces(w, r)
	case http.MethodPost:
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
		fmt.Println(err)
	}
	models.Places = append(models.Places, place)
	fmt.Println("post people")
}
