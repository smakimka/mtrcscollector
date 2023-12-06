package main

import (
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/smakimka/mtrcscollector/cmd/server/handlers"
	mw "github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func GetRouter() chi.Router {
	logger := log.New(os.Stdout, "", 3)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	storageMW := mw.WithMemStorage{S: s}

	r := chi.NewRouter()
	r.Use(chiMW.Logger)
	r.Use(storageMW.WithMemStorage)

	r.Get("/", handlers.GetAllMetrics)
	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Post("/{metricName}/{metricValue}", handlers.UpdateMetricHandler)
	})
	r.Route("/value/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Get("/{metricName}", handlers.GetMetricValueHandler)
	})

	return r
}
