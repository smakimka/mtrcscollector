package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

type UpdateMetricHandler struct {
	s storage.Storage
}

func NewUpdateMetricHandler(s storage.Storage) UpdateMetricHandler {
	return UpdateMetricHandler{s: s}
}

func (h UpdateMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch chi.URLParam(r, "metricKind") {
	case model.Gauge:
		value, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		err = h.s.UpdateGaugeMetric(model.GaugeMetric{
			Name:  chi.URLParam(r, "metricName"),
			Value: value,
		})

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "")

	case model.Counter:
		value, err := strconv.ParseInt(chi.URLParam(r, "metricValue"), 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		err = h.s.UpdateCounterMetric(model.CounterMetric{
			Name:  chi.URLParam(r, "metricName"),
			Value: value,
		})

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "")
	}
}
