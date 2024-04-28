package recommend

import (
	"fmt"
	"log/slog"
	"net/http"
	"places/internal/entities"
	response "places/internal/lib/api/response"
	"places/internal/storage"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Name   string           `json:"name"`
	Places []entities.Place `json:"places"`
}

func GetClosestPlaces(esStore *storage.ElasticStore, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.recommend.GetClosestPlaces"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		latitudeParam := r.URL.Query().Get("lat")
		longitudeParam := r.URL.Query().Get("lon")
		latitude, _ := strconv.ParseFloat(latitudeParam, 64)
		longitude, _ := strconv.ParseFloat(longitudeParam, 64)

		if latitude < 0 {
			render.JSON(w, r, response.Error(fmt.Sprintf("Invalid 'lat' value: '%v'", latitude)))
			return
		}

		if longitude < 0 {
			render.JSON(w, r, response.Error(fmt.Sprintf("Invalid 'lon' value: '%v'", longitude)))
			return
		}

		data, _, _ := esStore.GetClosestPlaces(longitude, latitude, logger)

		responseParams := Response{
			Name:   "Recommendation",
			Places: data,
		}

		responseOK(w, r, &responseParams)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, responseParams *Response) {
	render.JSON(w, r, Response{
		Name:   responseParams.Name,
		Places: responseParams.Places,
	})
}
