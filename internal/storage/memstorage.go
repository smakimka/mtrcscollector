package storage

import (
	"log"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
)

type MemStorage struct {
	Logger         *log.Logger
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func (s *MemStorage) Update(m mtrcs.Metric) error {
	var err error
	switch metric := m.(type) {
	case mtrcs.GaugeMetric:
		err = s.updateGaugeMetric(metric)
	case mtrcs.CounterMetric:
		err = s.updateCounterMetric(metric)
	}

	return err
}

func (s *MemStorage) Init() error {
	s.GaugeMetrics = make(map[string]float64)
	s.CounterMetrics = make(map[string]int64)
	return nil
}

func (s *MemStorage) updateGaugeMetric(m mtrcs.GaugeMetric) error {
	s.GaugeMetrics[m.Name] = m.Value
	s.Logger.Printf("updated gauge metric \"%s\" to %f", m.Name, m.Value)
	return nil
}

func (s *MemStorage) updateCounterMetric(m mtrcs.CounterMetric) error {
	s.CounterMetrics[m.Name] += m.Value
	s.Logger.Printf("updated counter metric \"%s\" to %d", m.Name, s.CounterMetrics[m.Name])
	return nil
}
