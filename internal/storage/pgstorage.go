package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smakimka/mtrcscollector/internal/model"
)

type PGStorage struct {
	p *pgxpool.Pool
}

func NewPGStorage(p *pgxpool.Pool) PGStorage {
	return PGStorage{
		p: p,
	}
}

func (s PGStorage) UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error) {
	return 0, nil
}

func (s PGStorage) UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	return nil
}
