package routes

import (
	"net/http"
	"places/internal/db"
	"places/internal/serializers"
	"text/template"
)

func GetPlacesHandler(w http.ResponseWriter, r *http.Request, esStore *db.ElasticStore) {
	tmpl := template.Must(template.ParseFiles("/app/internal/templates/index.gohtml"))
	data, _, _ := esStore.GetPlaces(10, 10)

	response := serializers.GetPlacesResponse{
		Total:  10,
		Places: data,
	}
	tmpl.Execute(w, response)
}
