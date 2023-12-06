package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
)

func MetricKind(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricKind := chi.URLParam(r, "metricKind")
		if metricKind != mtrcs.Gauge && metricKind != mtrcs.Counter {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("wrong metric type"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
