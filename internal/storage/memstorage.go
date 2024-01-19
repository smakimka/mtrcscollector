package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
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

type SaveData struct {
	GaugeMetrics   map[string]float64 `json:"gauge_metrics"`
	CounterMetrics map[string]int64   `json:"counter_metrics"`
}

func (s *MemStorage) Save(filePath string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data := SaveData{}

	data.GaugeMetrics = s.gaugeMetrics
	data.CounterMetrics = s.counterMetrics

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	os.WriteFile(filePath, byteData, fs.FileMode(0644))

	return nil
}

func (s *MemStorage) Restore(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	metricsData := SaveData{}
	if err = json.Unmarshal(data, &metricsData); err != nil {
		return err
	}

	s.gaugeMetrics = metricsData.GaugeMetrics
	s.counterMetrics = metricsData.CounterMetrics

	return nil
}

func (s *MemStorage) GetGaugeMetric(ctx context.Context, name string) (model.GaugeMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	gaugeVal, ok := s.gaugeMetrics[name]
	if ok {
		return model.GaugeMetric{Name: name, Value: gaugeVal}, nil
	}

	return model.GaugeMetric{}, ErrNoSuchMetric
}

func (s *MemStorage) GetCounterMetric(ctx context.Context, name string) (model.CounterMetric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	counterVal, ok := s.counterMetrics[name]
	if ok {
		return model.CounterMetric{Name: name, Value: counterVal}, nil
	}

	return model.CounterMetric{}, ErrNoSuchMetric
}

func (s *MemStorage) GetAllGaugeMetrics(ctx context.Context) ([]model.GaugeMetric, error) {
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

func (s *MemStorage) GetAllCounterMetrics(ctx context.Context) ([]model.CounterMetric, error) {
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

func (s *MemStorage) UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.gaugeMetrics[m.Name] = m.Value
	logger.Log.Debug().Msg(fmt.Sprintf("updated gauge metric \"%s\" to %f", m.Name, m.Value))

	return nil
}

func (s *MemStorage) UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.counterMetrics[m.Name] += m.Value
	logger.Log.Debug().Msg(fmt.Sprintf("updated counter metric \"%s\" to %d", m.Name, s.counterMetrics[m.Name]))
	return s.counterMetrics[m.Name], nil
}

func (s *MemStorage) UpdateMetrics(ctx context.Context, metricsData model.MetricsData) (model.MetricsData, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range metricsData {
		metricData := metricsData[i]

		switch metricData.Kind {
		case model.Gauge:
			s.gaugeMetrics[metricData.Name] = *metricData.Value
			logger.Log.Debug().Msg(fmt.Sprintf("updated gauge metric \"%s\" to %f", metricData.Name, *metricData.Value))
		case model.Counter:
			s.counterMetrics[metricData.Name] += *metricData.Delta
			newValue := s.counterMetrics[metricData.Name]
			metricsData[i].Delta = &newValue
			logger.Log.Debug().Msg(fmt.Sprintf("updated counter metric \"%s\" to %d", metricData.Name, newValue))
		}
	}

	return metricsData, nil
}
