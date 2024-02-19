package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

type GetMetricValueHandler struct {
	s storage.Storage
}

func NewGetMetricValueHandler(s storage.Storage) GetMetricValueHandler {
	return GetMetricValueHandler{s: s}
}

func (h GetMetricValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	switch chi.URLParam(r, "metricKind") {
	case model.Gauge:
		metric, err := h.s.GetGaugeMetric(ctx, chi.URLParam(r, "metricName"))

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
		metric, err := h.s.GetCounterMetric(ctx, chi.URLParam(r, "metricName"))

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

type ValueHandler struct {
	s storage.Storage
}

func NewValueHandler(s storage.Storage) ValueHandler {
	return ValueHandler{s: s}
}

func (h ValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	data := &model.MetricData{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
		return
	}

	switch data.Kind {
	case model.Gauge:
		m, err := h.s.GetGaugeMetric(ctx, data.Name)
		if err != nil {
			if err == storage.ErrNoSuchMetric {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
				return
			}

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
			return
		}

		data.Value = &m.Value
	case model.Counter:
		m, err := h.s.GetCounterMetric(ctx, data.Name)
		if err != nil {
			if err == storage.ErrNoSuchMetric {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
				return
			}

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
			return
		}

		data.Delta = &m.Value
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}
