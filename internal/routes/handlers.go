package routes

import (
	"fmt"
	"log"
	"net/http"
	"places/internal/db"
	"places/internal/entities"
	"strconv"
	"text/template"
)

type GetPlacesResponse struct {
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	Places   []entities.Place `json:"places"`
	PrevPage int              `json:"prev_page"`
	NextPage int              `json:"next_page"`
	LastPage int              `json:"last_page"`
}

func GetPlacesHandler(w http.ResponseWriter, r *http.Request, esStore *db.ElasticStore) {
	pageParam := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	log.Println("page", page)

	limit := 10
	offset := (page - 1) * limit
	log.Println("offset: ", offset)

	tmpl := template.Must(template.ParseFiles("/app/internal/templates/index.gohtml"))
	data, total, _ := esStore.GetPlaces(limit, offset)

	if page < 1 || page > total {
		w.WriteHeader(400)
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%v'", pageParam), http.StatusBadRequest)
		return
	}

	log.Println(data)
	response := GetPlacesResponse{
		Total:  total,
		Places: data,
	}
	if offset > 0 {
		response.PrevPage = page - 1
	}

	if offset+limit < total {
		response.NextPage = page + 1
	}

	response.LastPage = (total + limit - 1) / limit
	tmpl.Execute(w, response)
}
