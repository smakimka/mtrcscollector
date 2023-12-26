package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/smakimka/mtrcscollector/internal/server/handlers"
	"github.com/smakimka/mtrcscollector/internal/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func GetRouter() chi.Router {
	s := storage.NewMemStorage()

	getAllMetricsHandler := handlers.NewGetAllMetricsHandler(s)
	updateMetricHandler := handlers.NewUpdateMetricHandler(s)
	getMetricValueHandler := handlers.NewGetMetricValueHandler(s)
	updateHandler := handlers.NewUpdateHandler(s)
	valueHandler := handlers.NewValueHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", getAllMetricsHandler.ServeHTTP)

	r.Post("/update/", updateHandler.ServeHTTP)
	r.Post("/value/", valueHandler.ServeHTTP)

	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(middleware.MetricKind)
		r.Post("/{metricName}/{metricValue}", updateMetricHandler.ServeHTTP)
	})
	r.Route("/value/{metricKind}", func(r chi.Router) {
		r.Use(middleware.MetricKind)
		r.Get("/{metricName}", getMetricValueHandler.ServeHTTP)
	})

	return r
}
