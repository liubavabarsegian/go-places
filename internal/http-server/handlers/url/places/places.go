package places

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"places/internal/entities"
	response "places/internal/lib/api/response"
	"places/internal/storage"
	"text/template"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL  string `json:"url"`
	Page int    `json:"page,omitempty"`
}

type Response struct {
	response.Response
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	Places   []entities.Place `json:"places"`
	PrevPage int              `json:"prev_page"`
	NextPage int              `json:"next_page"`
	LastPage int              `json:"last_page"`
}

func New(esStore *storage.ElasticStore, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.root.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			logger.Error("Request body is empty")
			render.JSON(w, r, response.Error("Empty request"))
			return
		}
		if err != nil {
			logger.Error("Failed to decode request body", err)
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}

		page := req.Page
		// page, _ := strconv.Atoi(pageParam)
		logger.Info("Requested page №", page)

		limit := 10
		offset := (page - 1) * limit

		tmpl := template.Must(template.ParseFiles("/app/internal/templates/index.gohtml"))
		data, total, _ := esStore.GetPlaces(limit, offset)

		if page < 1 || page > total {
			w.WriteHeader(400)
			http.Error(w, fmt.Sprintf("Invalid 'page' value: '%v'", page), http.StatusBadRequest)
			return
		}

		response := Response{
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
}
