package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	var m mtrcs.Metric
	var convErr error
	metricKind := chi.URLParam(r, "metricKind")

	switch metricKind {
	case mtrcs.Gauge:
		value, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64)
		convErr = err
		m = mtrcs.GaugeMetric{
			Name:  chi.URLParam(r, "metricName"),
			Value: value,
		}
	case mtrcs.Counter:
		value, err := strconv.ParseInt(chi.URLParam(r, "metricValue"), 10, 64)
		convErr = err
		m = mtrcs.CounterMetric{
			Name:  chi.URLParam(r, "metricName"),
			Value: value,
		}
	}

	if convErr != nil {
		http.Error(w, convErr.Error(), http.StatusBadRequest)
		return
	}

	s := r.Context().Value(middleware.StorageKey).(storage.Storage)
	err := s.UpdateMetric(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
