package storage

import (
	"fmt"
	"sync"

	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
)

type MemStorage struct {
	mutex          sync.RWMutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	return s
}

func (s *MemStorage) GetGaugeMetric(name string) (model.GaugeMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	gaugeVal, ok := s.gaugeMetrics[name]
	if ok {
		return model.GaugeMetric{Name: name, Value: gaugeVal}, nil
	}

	return model.GaugeMetric{}, ErrNoSuchMetric
}

func (s *MemStorage) GetCounterMetric(name string) (model.CounterMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	counterVal, ok := s.counterMetrics[name]
	if ok {
		return model.CounterMetric{Name: name, Value: counterVal}, nil
	}

	return model.CounterMetric{}, ErrNoSuchMetric
}

func (s *MemStorage) GetAllGaugeMetrics() ([]model.GaugeMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	metrics := make([]model.GaugeMetric, len(s.gaugeMetrics))
	idx := 0

	for name, value := range s.gaugeMetrics {
		metrics[idx] = model.GaugeMetric{Name: name, Value: value}
		idx++
	}

	return metrics, nil
}

func (s *MemStorage) GetAllCounterMetrics() ([]model.CounterMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	metrics := make([]model.CounterMetric, len(s.counterMetrics))
	idx := 0

	for name, value := range s.counterMetrics {
		metrics[idx] = model.CounterMetric{Name: name, Value: value}
		idx++
	}

	return metrics, nil
}

func (s *MemStorage) UpdateGaugeMetric(m model.GaugeMetric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.gaugeMetrics[m.Name] = m.Value
	logger.Log.Debug().Msg(fmt.Sprintf("updated gauge metric \"%s\" to %f", m.Name, m.Value))

	return nil
}

func (s *MemStorage) UpdateCounterMetric(m model.CounterMetric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.counterMetrics[m.Name] += m.Value
	logger.Log.Debug().Msg(fmt.Sprintf("updated counter metric \"%s\" to %d", m.Name, s.counterMetrics[m.Name]))
	return nil
}
