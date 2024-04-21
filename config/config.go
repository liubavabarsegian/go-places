package config

import (
	"fmt"
	"net/http"
	"places/internal/routes"

	"github.com/gorilla/mux"
)

const PlacesFilePath = "/app/config/data.csv"

func ConfigServer() {
	r := mux.NewRouter()
	routes.RegisterPlacesRoutes(r)
	http.Handle("/", r)
	fmt.Println("config")
	http.ListenAndServe("localhost:8888", r)
	fmt.Println("config")
}
