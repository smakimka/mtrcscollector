package storage

import (
	"errors"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
)

var (
	ErrNoSuchMetric = errors.New("no such metric")
)

type Storage interface {
	Init() error
	UpdateMetric(m mtrcs.Metric) error
	GetMetric(kind string, name string) (mtrcs.Metric, error)
	GetAllMetrics() ([]mtrcs.Metric, error)
}
