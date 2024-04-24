package routes

import (
	"net/http"
	"places/internal/entities"
	"places/internal/serializers"
	"text/template"
)

func GetPlacesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/app/internal/templates/index.gohtml"))
	data := serializers.PlacePageData{
		Places: []entities.Place{
			{ID: 1, Name: "aaA", Address: "kargina", Phone: "+7"},
			{ID: 2, Name: "aaa"},
			{ID: 3, Name: "true"},
		},
	}
	tmpl.Execute(w, data)
}
