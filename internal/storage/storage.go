package storage

import (
	"context"
	"errors"

	"github.com/smakimka/mtrcscollector/internal/model"
)

var (
	ErrNoSuchMetric = errors.New("no such metric")
)

type Updater interface {
	UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error)
	UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error
}

type Getter interface {
	GetGaugeMetric(ctx context.Context, name string) (model.GaugeMetric, error)
	GetCounterMetric(ctx context.Context, name string) (model.CounterMetric, error)
	GetAllGaugeMetrics(ctx context.Context) ([]model.GaugeMetric, error)
	GetAllCounterMetrics(ctx context.Context) ([]model.CounterMetric, error)
}

type Storage interface {
	Updater
	Getter

	Restore(filePath string) error
	Save(filePath string) error
}
