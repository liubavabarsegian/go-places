package places

import (
	"fmt"
	"log/slog"
	"net/http"
	"places/internal/config"
	"places/internal/entities"
	response "places/internal/lib/api/response"
	"places/internal/storage"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	Places   []entities.Place `json:"places"`
	PrevPage int              `json:"prev_page"`
	NextPage int              `json:"next_page"`
	LastPage int              `json:"last_page"`
}

func GetPlaces(esStore *storage.ElasticStore, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.places.GetPlaces"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pageParam := r.URL.Query().Get("page")
		page, _ := strconv.Atoi(pageParam)
		logger.Info("Requested page №", pageParam)

		limit := 10
		offset := (page - 1) * limit
		data, total, _ := esStore.GetPlaces(limit, offset, logger)

		if page < 1 || page > total {
			render.JSON(w, r, response.Error(fmt.Sprintf("Invalid 'page' value: '%v'", page)))
			return
		}

		responseParams := Response{
			Total:  total,
			Places: data,
		}
		if offset > 0 {
			responseParams.PrevPage = page - 1
		}
		if offset+limit < total {
			responseParams.NextPage = page + 1
		}
		responseParams.LastPage = (total + limit - 1) / limit

		responseOK(w, r, &responseParams)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, responseParams *Response) {
	render.JSON(w, r, Response{
		Total:    responseParams.Total,
		Name:     config.IndexName,
		Places:   responseParams.Places,
		PrevPage: responseParams.PrevPage,
		NextPage: responseParams.NextPage,
		LastPage: responseParams.LastPage,
	})
}
