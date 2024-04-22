package routes

import (
	"fmt"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	RegisterPlacesRoutes(mux)
	return mux
}

func RegisterPlacesRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", GetPlacesHandler)
	fmt.Println("Registered routes")
}
