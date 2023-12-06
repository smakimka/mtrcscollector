package main

import (
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	tests := []struct {
		name               string
		callTimes          int
		wantGaugeLength    int
		wantCounterLength  int
		wantPollCountValue int
	}{
		{
			name:               "single update",
			callTimes:          1,
			wantGaugeLength:    26,
			wantCounterLength:  1,
			wantPollCountValue: 1,
		},
		{
			name:               "double updarte",
			callTimes:          2,
			wantGaugeLength:    26,
			wantCounterLength:  1,
			wantPollCountValue: 2,
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &storage.MemStorage{Logger: logger}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := s.Init()
			require.NoError(t, err, "Error initializing memstorage")
			for i := 0; i < test.callTimes; i++ {
				m := runtime.MemStats{}
				runtime.ReadMemStats(&m)
				updateMetrics(&m, s, logger)

				assert.Equal(t, test.wantGaugeLength, len(s.GaugeMetrics))
				assert.Equal(t, test.wantCounterLength, len(s.CounterMetrics))
			}

			pollCount, err := s.GetMetric("counter", "PollCount")
			require.NoError(t, err)
			assert.Equal(t, int64(test.wantPollCountValue), pollCount.GetValue())
		})
	}
}
