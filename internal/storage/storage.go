package storage

import (
	"errors"

	"github.com/smakimka/mtrcscollector/internal/model"
)

var (
	ErrNoSuchMetric = errors.New("no such metric")
)

type Storage interface {
	UpdateCounterMetric(m model.CounterMetric) (int64, error)
	UpdateGaugeMetric(m model.GaugeMetric) error

	GetGaugeMetric(name string) (model.GaugeMetric, error)
	GetCounterMetric(name string) (model.CounterMetric, error)

	GetAllGaugeMetrics() ([]model.GaugeMetric, error)
	GetAllCounterMetrics() ([]model.CounterMetric, error)

	Restore(filePath string) error
}
