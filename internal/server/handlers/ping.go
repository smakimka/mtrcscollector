package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PingHandler struct {
	p *pgxpool.Pool
}

func NewPingHandler(pool *pgxpool.Pool) PingHandler {
	return PingHandler{
		p: pool,
	}
}

func (h PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := h.p.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
