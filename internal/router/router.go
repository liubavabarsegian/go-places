package router

import (
	"log/slog"
	"places/internal/http-server/handlers/url/places"
	"places/internal/http-server/handlers/url/root"
	"places/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func RegisterPlacesRoutes(esStore *storage.ElasticStore, router *chi.Mux, logger *slog.Logger) {
	router.Get("/", root.New(esStore, logger))
	router.Get("/api/places", places.New(esStore, logger))
	logger.Info("Registered routes")
}

func SetUpRouter(esStore *storage.ElasticStore, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	RegisterPlacesRoutes(esStore, router, logger)
	return router
}
