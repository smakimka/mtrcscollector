package storage

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/logger"
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

	logger.SetLevel(logger.Debug)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.UpdateGaugeMetric(ctx, test.newMetric)
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

	logger.SetLevel(logger.Debug)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.counterMetrics = test.counterMetrics
			s.UpdateCounterMetric(ctx, test.newMetric)
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

	logger.SetLevel(logger.Debug)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.counterMetrics = test.counterMetrics

			switch test.metricKind {
			case model.Gauge:
				m, err := s.GetGaugeMetric(ctx, test.metricName)
				if test.wantErr {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, test.wantMetric, m)
				}
			case model.Counter:
				m, err := s.GetCounterMetric(ctx, test.metricName)
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

	logger.SetLevel(logger.Debug)
	s := &MemStorage{
		mutex:          sync.RWMutex{},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.gaugeMetrics = test.gaugeMetrics
			s.counterMetrics = test.counterMetrics

			gaugeMetrics, err := s.GetAllGaugeMetrics(ctx)
			require.NoError(t, err, "error getting all gauge metrics from memstorage")
			assert.ElementsMatch(t, test.wantGaugeMetrics, gaugeMetrics)

			counterMetrics, err := s.GetAllCounterMetrics(ctx)
			require.NoError(t, err, "error getting all gauge metrics from memstorage")
			assert.ElementsMatch(t, test.wantCounterMetrics, counterMetrics)
		})
	}
}

func TestSaveLoad(t *testing.T) {
	type Data struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name string
		save Data
		load Data
	}{
		{
			name: "load save test #1",
			save: Data{
				gaugeMetrics:   map[string]float64{"test_1": 1.1, "test_2": 2.2},
				counterMetrics: map[string]int64{"test_1": 1, "test_2": 2},
			},
			load: Data{
				map[string]float64{"test_1": 1.1, "test_2": 2.2},
				map[string]int64{"test_1": 1, "test_2": 2},
			},
		},
	}

	testFilePath := fmt.Sprintf("/tmp/test%d.json", rand.Int63())
	defer os.Remove(testFilePath)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &MemStorage{
				mutex:          sync.RWMutex{},
				gaugeMetrics:   test.save.gaugeMetrics,
				counterMetrics: test.save.counterMetrics,
			}

			err := s.Save(testFilePath)
			require.NoError(t, err)

			err = s.Restore(testFilePath)
			require.NoError(t, err)

			assert.Equal(t, test.load.gaugeMetrics, s.gaugeMetrics)
			assert.Equal(t, test.load.counterMetrics, s.counterMetrics)
		})
	}
}
