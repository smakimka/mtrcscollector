// Модуль storage осуществляет хранение метрик на агенте и на сервере
package storage

import (
	"context"
	"errors"

	"github.com/smakimka/mtrcscollector/internal/model"
)

var (
	ErrNoSuchMetric             = errors.New("no such metric")
	_               Storage     = (*MemStorage)(nil)
	_               Storage     = (*PGStorage)(nil)
	_               SyncStorage = (*SyncMemStorage)(nil)
)

type updater interface {
	UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error)
	UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error
	UpdateMetrics(ctx context.Context, metricsData model.MetricsData) error
}

type getter interface {
	GetGaugeMetric(ctx context.Context, name string) (model.GaugeMetric, error)
	GetCounterMetric(ctx context.Context, name string) (model.CounterMetric, error)
	GetAllGaugeMetrics(ctx context.Context) ([]model.GaugeMetric, error)
	GetAllCounterMetrics(ctx context.Context) ([]model.CounterMetric, error)
}

// Основной интерфейс, который реализуют все хранилища
type Storage interface {
	updater
	getter
}

// Интерфейс для хранилищ, которым нужно переодически сохранять данные и потом восстанавливаться из сохранения
type SyncStorage interface {
	Storage
	Restore(filePath string) error
	Save(filePath string) error
}
