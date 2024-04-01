package handlers

import (
	"context"
	"net/http"

	"github.com/smakimka/mtrcscollector/internal/storage"
)

type PingHandler struct {
	s storage.Storage
}

func NewPingHandler(s storage.Storage) PingHandler {
	return PingHandler{
		s: s,
	}
}

// Ping godoc
// @Tags Status
// @Summary Запрос для проверки соединения с БД
// @ID Ping
// @Accept  plain
// @Produce plain
// @Success 200 {string} string ""
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /ping [get]
func (h PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	switch s := h.s.(type) {
	case storage.PGStorage:
		err := s.Ping(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
