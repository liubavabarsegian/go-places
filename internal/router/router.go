package router

import (
	"log/slog"
	"places/internal/http-server/handlers/url/auth"
	"places/internal/http-server/handlers/url/places"
	"places/internal/http-server/handlers/url/recommend"
	"places/internal/http-server/handlers/url/root"
	jwt_middleware "places/internal/http-server/middleware/jwt_middleware"

	"places/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func RegisterPlacesRoutes(esStore *storage.ElasticStore, router *chi.Mux, logger *slog.Logger) {
	router.Get("/", root.GetPlaces(esStore, logger))
	router.Get("/api/places", places.GetPlaces(esStore, logger))
	router.Get("/api/get_token", auth.GetToken(logger))
	router.Get("/api/recommend", jwt_middleware.JWTMiddleware(recommend.GetClosestPlaces(esStore, logger)))

	logger.Info("Registered routes")
}

func SetUpRouter(esStore *storage.ElasticStore, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	// CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов
	router.Use(cors.Handler)

	RegisterPlacesRoutes(esStore, router, logger)
	return router
}
