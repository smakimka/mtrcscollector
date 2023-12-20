package router

import (
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"

	"github.com/smakimka/mtrcscollector/internal/server/handlers"
	mw "github.com/smakimka/mtrcscollector/internal/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func GetRouter() chi.Router {
	l := log.New(os.Stdout, "", 3)

	s := storage.NewMemStorage().WithLogger(l)

	getAllMetricsHandler := handlers.NewGetAllMetricsHandler().WithStorage(s).WithLogger(l)
	updateMetricHandler := handlers.NewUpdateMetricHandler().WithStorage(s).WithLogger(l)
	getMetricValueHandler := handlers.NewGetMetricValueHandler().WithStorage(s).WithLogger(l)

	r := chi.NewRouter()
	r.Use(chiMW.Logger)

	r.Get("/", getAllMetricsHandler.ServeHTTP)
	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Post("/{metricName}/{metricValue}", updateMetricHandler.ServeHTTP)
	})
	r.Route("/value/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Get("/{metricName}", getMetricValueHandler.ServeHTTP)
	})

	return r
}
