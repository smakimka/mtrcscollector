package handlers

import (
	"context"
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

// UpdateOld godoc
// @Tags Update
// @Summary Запрос для обновления метрики
// @ID UpdateOld
// @Accept  plain
// @Produce plain
// @Param metricKind path string true "Тип метрики для обновления"
// @Param metricName path string true "имя метрики"
// @Param metricValue path string true "Значение метрики"
// @Success 200 {string} string "20"
// @Failure 400 {string} string "ошибка"
// @Failure 500 {object} string "ошибка"
// @Router /update/{metricKind}/{metricName}/{metricValue} [post]
func (h UpdateMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	switch chi.URLParam(r, "metricKind") {
	case model.Gauge:
		value, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		err = h.s.UpdateGaugeMetric(ctx, model.GaugeMetric{
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

		_, err = h.s.UpdateCounterMetric(ctx, model.CounterMetric{
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

type UpdateHandler struct {
	s storage.Storage
}

func NewUpdateHandler(s storage.Storage) UpdateHandler {
	return UpdateHandler{s: s}
}

// Update godoc
// @Tags Update
// @Summary Запрос для обновления метрики
// @ID Update
// @Accept  json
// @Produce json
// @Param metric body model.MetricData true "Метрика для обновления"
// @Success 200 {object} model.MetricData
// @Failure 500 {object} model.Response
// @Router /update/ [post]
func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		if data.Value == nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, model.Response{Ok: false, Detail: model.ErrMissingFields.Error()})
			return
		}

		err := h.s.UpdateGaugeMetric(ctx, model.GaugeMetric{
			Name:  data.Name,
			Value: *data.Value,
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
			return
		}
	case model.Counter:
		if data.Delta == nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, model.Response{Ok: false, Detail: model.ErrMissingFields.Error()})
			return
		}

		newVal, err := h.s.UpdateCounterMetric(ctx, model.CounterMetric{
			Name:  data.Name,
			Value: *data.Delta,
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
			return
		}
		data.Delta = &newVal
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}
