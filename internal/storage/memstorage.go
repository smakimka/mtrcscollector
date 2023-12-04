package storage

import (
	"log"
	"sync"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
)

type MemStorage struct {
	Logger         *log.Logger
	Mutex          sync.RWMutex
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func (s *MemStorage) UpdateMetric(m mtrcs.Metric) error {
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
	s.Mutex = sync.RWMutex{}
	s.GaugeMetrics = make(map[string]float64)
	s.CounterMetrics = make(map[string]int64)
	return nil
}

func (s *MemStorage) GetMetric(name string) (mtrcs.Metric, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	gaugeVal, ok := s.GaugeMetrics[name]
	if ok {
		return mtrcs.GaugeMetric{Name: name, Value: gaugeVal}, nil
	}

	counterVal, ok := s.CounterMetrics[name]
	if ok {
		return mtrcs.CounterMetric{Name: name, Value: counterVal}, nil
	}

	return nil, ErrNoSuchMetric
}

func (s *MemStorage) GetAllMetrics() ([]mtrcs.Metric, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	metrics := make([]mtrcs.Metric, len(s.CounterMetrics)+len(s.GaugeMetrics))
	idx := 0

	for name, value := range s.GaugeMetrics {
		metrics[idx] = mtrcs.GaugeMetric{Name: name, Value: value}
		idx++
	}

	for name, value := range s.CounterMetrics {
		metrics[idx] = mtrcs.CounterMetric{Name: name, Value: value}
		idx++
	}

	return metrics, nil
}

func (s *MemStorage) updateGaugeMetric(m mtrcs.GaugeMetric) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.GaugeMetrics[m.Name] = m.Value
	s.Logger.Printf("updated gauge metric \"%s\" to %f", m.Name, m.Value)

	return nil
}

func (s *MemStorage) updateCounterMetric(m mtrcs.CounterMetric) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.CounterMetrics[m.Name] += m.Value
	s.Logger.Printf("updated counter metric \"%s\" to %d", m.Name, s.CounterMetrics[m.Name])

	return nil
}
