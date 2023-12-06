package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func GetMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value(middleware.StorageKey).(storage.Storage)
	metric, err := s.GetMetric(chi.URLParam(r, "metricKind"), chi.URLParam(r, "metricName"))

	if err != nil {
		if err == storage.ErrNoSuchMetric {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.GetStringValue()))
}
