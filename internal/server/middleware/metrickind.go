package middleware

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smakimka/mtrcscollector/internal/model"
)

var ErrWrongMetricKind = errors.New("wrong metric type")

func MetricKind(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricKind := chi.URLParam(r, "metricKind")
		if metricKind != model.Gauge && metricKind != model.Counter {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(ErrWrongMetricKind.Error()))
			return
		}

		next.ServeHTTP(w, r)
	})
}
