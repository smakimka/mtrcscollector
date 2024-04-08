package router

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/smakimka/mtrcscollector/internal/server/handlers"
	"github.com/smakimka/mtrcscollector/internal/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

// @Title mtrcscollector API
// @Description Серви для сбора метрик.
// @Version 1.0

// @BasePath /

// @Tag.name Get
// @Tag.description "Группа запросов получения метрик"

// @Tag.name Update
// @Tag.description "Группа запросов обновления метрик"

// @Tag.name Status
// @Tag.description "Группа запросов статуса сервиса"

func GetRouter(s storage.Storage) chi.Router {
	getAllMetricsHandler := handlers.NewGetAllMetricsHandler(s)
	updateMetricHandler := handlers.NewUpdateMetricHandler(s)
	getMetricValueHandler := handlers.NewGetMetricValueHandler(s)
	updateHandler := handlers.NewUpdateHandler(s)
	valueHandler := handlers.NewValueHandler(s)
	pingHandler := handlers.NewPingHandler(s)
	updatesHandler := handlers.NewUpdatesHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Auth)
	r.Use(middleware.Gzip)

	r.Get("/ping", pingHandler.ServeHTTP)

	r.Route("/", func(r chi.Router) {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline/", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile/", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol/", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace/", pprof.Trace)

		r.Handle("/debug/pprof/allocs/", pprof.Handler("allocs"))
		r.Handle("/debug/pprof/block/", pprof.Handler("block"))
		r.Handle("/debug/pprof/goroutine/", pprof.Handler("goroutine"))
		r.Handle("/debug/pprof/heap/", pprof.Handler("heap"))
		r.Handle("/debug/pprof/mutex/", pprof.Handler("mutex"))
		r.Handle("/debug/pprof/threadcreate/", pprof.Handler("threadcreate"))

		r.Get("/", getAllMetricsHandler.ServeHTTP)
		r.Post("/update/", updateHandler.ServeHTTP)
		r.Post("/updates/", updatesHandler.ServeHTTP)
		r.Post("/value/", valueHandler.ServeHTTP)
	})

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
