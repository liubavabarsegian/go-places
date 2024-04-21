package routes

import (
	"fmt"

	"github.com/gorilla/mux"
)

var RegisterPlacesRoutes = func(router *mux.Router) {

	router.HandleFunc("/places", PlacesHandler) //.Methods("POST")
	fmt.Println("Configured routes")
}
