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

// GetMetricValue godoc
// @Tags Get
// @Summary Запрос для получения метрики
// @ID GetMetricValue
// @Accept  plain
// @Produce plain
// @Param metricKind path string true "Тип метрики для обновления"
// @Param metricName path string true "имя метрики"
// @Success 200 {string} string "20"
// @Failure 500 {string} string "ошибка"
// @Failure 404 {string} string ""
// @Router /value/{metricKind}/{metricName} [get]
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

// Value godoc
// @Tags Get
// @Summary Запрос для получения метрики
// @ID Value
// @Accept  json
// @Produce json
// @Param metric body model.MetricData true "Метрика для обновления"
// @Success 200 {object} model.MetricData
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /value/ [get]
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
