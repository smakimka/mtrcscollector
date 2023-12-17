package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

type GetMetricValueHandler struct {
	s storage.Storage
	l *log.Logger
}

func NewGetMetricValueHandler() GetMetricValueHandler {
	return GetMetricValueHandler{}
}

func (h GetMetricValueHandler) WithLogger(l *log.Logger) GetMetricValueHandler {
	h.l = l
	return h
}

func (h GetMetricValueHandler) WithStorage(s storage.Storage) GetMetricValueHandler {
	h.s = s
	return h
}

func (h GetMetricValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch chi.URLParam(r, "metricKind") {
	case model.Gauge:
		metric, err := h.s.GetGaugeMetric(chi.URLParam(r, "metricName"))

		if err != nil {
			if err == storage.ErrNoSuchMetric {
				render.Status(r, http.StatusNotFound)
				render.PlainText(w, r, err.Error())
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, metric.GetStringValue())

	case model.Counter:
		metric, err := h.s.GetCounterMetric(chi.URLParam(r, "metricName"))

		if err != nil {
			if err == storage.ErrNoSuchMetric {
				render.Status(r, http.StatusNotFound)
				render.PlainText(w, r, err.Error())
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, metric.GetStringValue())
	}

}
