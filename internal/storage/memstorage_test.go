package storage

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/model"
)

func TestUpdateGaugeMetric(t *testing.T) {
	tests := []struct {
		name         string
		gaugeMetrics map[string]float64
		newMetric    model.GaugeMetric
		want         map[string]float64
	}{
		{
			name:         "create new metric",
			gaugeMetrics: map[string]float64{},
			newMetric:    model.GaugeMetric{Name: "test", Value: 1.0},
			want:         map[string]float64{"test": 1.0},
		},
		{
			name:         "update metric",
			gaugeMetrics: map[string]float64{"test": 1.0},
			newMetric:    model.GaugeMetric{Name: "test", Value: 5.0},
			want:         map[string]float64{"test": 5.0},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		logger:         logger,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.UpdateGaugeMetric(test.newMetric)
			assert.Equal(t, test.want, s.gaugeMetrics)
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	tests := []struct {
		name           string
		counterMetrics map[string]int64
		newMetric      model.CounterMetric
		want           map[string]int64
	}{
		{
			name:           "create new metric",
			counterMetrics: map[string]int64{},
			newMetric:      model.CounterMetric{Name: "test", Value: 1},
			want:           map[string]int64{"test": 1},
		},
		{
			name:           "update metric",
			counterMetrics: map[string]int64{"test": 1},
			newMetric:      model.CounterMetric{Name: "test", Value: 5},
			want:           map[string]int64{"test": 6},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		logger:         logger,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.counterMetrics = test.counterMetrics
			s.UpdateCounterMetric(test.newMetric)
			assert.Equal(t, test.want, s.counterMetrics)
		})
	}
}

func TestGetMetric(t *testing.T) {
	tests := []struct {
		name           string
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
		metricKind     string
		metricName     string
		wantErr        bool
		wantMetric     interface{}
	}{
		{
			name:           "get gauge metric",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{},
			metricKind:     "gauge",
			metricName:     "test",
			wantErr:        false,
			wantMetric:     model.GaugeMetric{Name: "test", Value: 1.0},
		},
		{
			name:           "get counter metric",
			gaugeMetrics:   map[string]float64{},
			counterMetrics: map[string]int64{"test": 1},
			metricKind:     "counter",
			metricName:     "test",
			wantErr:        false,
			wantMetric:     model.CounterMetric{Name: "test", Value: 1},
		},
		{
			name:           "get non existent gauge metric",
			gaugeMetrics:   map[string]float64{},
			counterMetrics: map[string]int64{},
			metricKind:     "gauge",
			metricName:     "test",
			wantErr:        true,
			wantMetric:     nil,
		},
		{
			name:           "get non existent counter metric",
			gaugeMetrics:   map[string]float64{},
			counterMetrics: map[string]int64{},
			metricKind:     "counter",
			metricName:     "test",
			wantErr:        true,
			wantMetric:     nil,
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		logger:         logger,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.counterMetrics = test.counterMetrics

			switch test.metricKind {
			case model.Gauge:
				m, err := s.GetGaugeMetric(test.metricName)
				if test.wantErr {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, test.wantMetric, m)
				}
			case model.Counter:
				m, err := s.GetCounterMetric(test.metricName)
				if test.wantErr {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, test.wantMetric, m)
				}
			}
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	tests := []struct {
		name               string
		gaugeMetrics       map[string]float64
		counterMetrics     map[string]int64
		wantGaugeMetrics   []model.GaugeMetric
		wantCounterMetrics []model.CounterMetric
	}{
		{
			name:               "empty storage",
			gaugeMetrics:       map[string]float64{},
			counterMetrics:     map[string]int64{},
			wantGaugeMetrics:   []model.GaugeMetric{},
			wantCounterMetrics: []model.CounterMetric{},
		},
		{
			name:           "non empty storage",
			gaugeMetrics:   map[string]float64{"test_1": 1.1, "test_2": 2.2},
			counterMetrics: map[string]int64{"test_1": 1, "test_2": 2},
			wantGaugeMetrics: []model.GaugeMetric{
				{Name: "test_1", Value: 1.1},
				{Name: "test_2", Value: 2.2},
			},
			wantCounterMetrics: []model.CounterMetric{
				{Name: "test_1", Value: 1},
				{Name: "test_2", Value: 2},
			},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		logger:         logger,
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.counterMetrics = test.counterMetrics

			gaugeMetrics, err := s.GetAllGaugeMetrics()
			require.NoError(t, err, "error getting all gauge metrics from memstorage")
			assert.ElementsMatch(t, test.wantGaugeMetrics, gaugeMetrics)

			counterMetrics, err := s.GetAllCounterMetrics()
			require.NoError(t, err, "error getting all gauge metrics from memstorage")
			assert.ElementsMatch(t, test.wantCounterMetrics, counterMetrics)
		})
	}
}
