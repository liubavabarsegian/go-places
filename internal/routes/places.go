package routes

import (
	"github.com/gorilla/mux"
)

var RegisterPlacesAppRoutes = func(router *mux.Router) {
	router.HandleFunc("/places", PlacesHandler) //.Methods("POST")
}
