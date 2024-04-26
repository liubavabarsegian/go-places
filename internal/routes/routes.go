package routes

import (
	"fmt"
	"net/http"
	"places/internal/db"
)

func RegisterRoutes(esStore *db.ElasticStore) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterPlacesRoutes(mux, esStore)
	return mux
}

func RegisterPlacesRoutes(mux *http.ServeMux, esStore *db.ElasticStore) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		GetPlacesHandler(w, r, esStore)
	})

	fmt.Println("Registered routes")
}
