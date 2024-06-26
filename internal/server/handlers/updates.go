package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/render"

	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

type UpdatesHandler struct {
	s storage.Storage
}

func NewUpdatesHandler(s storage.Storage) UpdatesHandler {
	return UpdatesHandler{s: s}
}

// Updates godoc
// @Tags Update
// @Summary Запрос для обновления метрики
// @ID Updates
// @Accept  json
// @Produce json
// @Param metric body model.MetricsData true "Метрики для обновления"
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /updates/ [post]
func (h UpdatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	data := new(model.MetricsData)
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
		return
	}

	err := h.s.UpdateMetrics(ctx, *data)
	if err != nil {
		logger.Log.Err(err).Msg("error updating metrics")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.Response{Ok: false, Detail: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.Response{Ok: true})
}
