package routes

import (
	"fmt"

	"github.com/gorilla/mux"
)

var RegisterPlacesAppRoutes = func(router *mux.Router) {
	router.HandleFunc("/places", PlacesHandler) //.Methods("POST")
	fmt.Println("AAAA")
}
