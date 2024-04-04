package main

import (
	"fmt"
	"log"
	"net/http"

	"PlacesApp/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("HELLO")
	r := mux.NewRouter()
	routes.RegisterPlacesAppRoutes(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9201", r))
	log.Println("listening")
}
