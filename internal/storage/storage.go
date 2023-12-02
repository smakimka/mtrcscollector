package storage

import "github.com/smakimka/mtrcscollector/internal/mtrcs"

type Storage interface {
	Init() error
	Update(m mtrcs.Metric) error
}
