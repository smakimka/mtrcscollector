package storage

import (
	"log"
	"os"
	"testing"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateGaugeMetric(t *testing.T) {
	tests := []struct {
		name         string
		gaugeMetrics map[string]float64
		newMetric    mtrcs.GaugeMetric
		want         map[string]float64
	}{
		{
			name:         "create new metric",
			gaugeMetrics: map[string]float64{},
			newMetric:    mtrcs.GaugeMetric{Name: "test", Value: 1.0},
			want:         map[string]float64{"test": 1.0},
		},
		{
			name:         "update metric",
			gaugeMetrics: map[string]float64{"test": 1.0},
			newMetric:    mtrcs.GaugeMetric{Name: "test", Value: 5.0},
			want:         map[string]float64{"test": 5.0},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{Logger: logger}
	err := s.Init()
	require.NoError(t, err, "error initializing memstorage")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.GaugeMetrics = test.gaugeMetrics
			s.updateGaugeMetric(test.newMetric)
			assert.Equal(t, test.want, s.GaugeMetrics)
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	tests := []struct {
		name           string
		counterMetrics map[string]int64
		newMetric      mtrcs.CounterMetric
		want           map[string]int64
	}{
		{
			name:           "create new metric",
			counterMetrics: map[string]int64{},
			newMetric:      mtrcs.CounterMetric{Name: "test", Value: 1},
			want:           map[string]int64{"test": 1},
		},
		{
			name:           "update metric",
			counterMetrics: map[string]int64{"test": 1},
			newMetric:      mtrcs.CounterMetric{Name: "test", Value: 5},
			want:           map[string]int64{"test": 6},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{Logger: logger}
	err := s.Init()
	require.NoError(t, err, "error initializing memstorage")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.CounterMetrics = test.counterMetrics
			s.updateCounterMetric(test.newMetric)
			assert.Equal(t, test.want, s.CounterMetrics)
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
		wantMetric     mtrcs.Metric
	}{
		{
			name:           "get gauge metric",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{},
			metricKind:     "gauge",
			metricName:     "test",
			wantErr:        false,
			wantMetric:     mtrcs.GaugeMetric{Name: "test", Value: 1.0},
		},
		{
			name:           "get counter metric",
			gaugeMetrics:   map[string]float64{},
			counterMetrics: map[string]int64{"test": 1},
			metricKind:     "counter",
			metricName:     "test",
			wantErr:        false,
			wantMetric:     mtrcs.CounterMetric{Name: "test", Value: 1},
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
	s := &MemStorage{Logger: logger}
	err := s.Init()
	require.NoError(t, err, "error initializing memstorage")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.GaugeMetrics = test.gaugeMetrics
			s.CounterMetrics = test.counterMetrics

			m, err := s.GetMetric(test.metricKind, test.metricName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.wantMetric, m)
			}
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	tests := []struct {
		name           string
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
		want           []mtrcs.Metric
	}{
		{
			name:           "empty storage",
			gaugeMetrics:   map[string]float64{},
			counterMetrics: map[string]int64{},
			want:           []mtrcs.Metric{},
		},
		{
			name:           "non empty storage",
			gaugeMetrics:   map[string]float64{"test": 1.0},
			counterMetrics: map[string]int64{"test": 1},
			want: []mtrcs.Metric{
				mtrcs.GaugeMetric{Name: "test", Value: 1.0},
				mtrcs.CounterMetric{Name: "test", Value: 1},
			},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &MemStorage{Logger: logger}
	err := s.Init()
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.GaugeMetrics = test.gaugeMetrics
			s.CounterMetrics = test.counterMetrics

			metrics, err := s.GetAllMetrics()
			require.NoError(t, err, "error getting all metrics from memstorage")
			assert.Equal(t, test.want, metrics)
		})
	}
}
