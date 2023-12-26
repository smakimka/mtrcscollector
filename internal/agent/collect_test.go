package agent

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/storage"
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
			name:               "double update",
			callTimes:          2,
			wantGaugeLength:    26,
			wantCounterLength:  1,
			wantPollCountValue: 2,
		},
	}

	logger.SetLevel(logger.Debug)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := storage.NewMemStorage()

			for i := 0; i < test.callTimes; i++ {
				m := runtime.MemStats{}
				runtime.ReadMemStats(&m)
				UpdateMetrics(&m, s)

				gaugeMetrics, err := s.GetAllGaugeMetrics()
				assert.NoError(t, err)
				counterMetrics, err := s.GetAllCounterMetrics()
				assert.NoError(t, err)
				assert.Equal(t, test.wantGaugeLength, len(gaugeMetrics))
				assert.Equal(t, test.wantCounterLength, len(counterMetrics))
			}

			pollCount, err := s.GetCounterMetric("PollCount")
			require.NoError(t, err)
			assert.Equal(t, int64(test.wantPollCountValue), pollCount.GetValue())
		})
	}
}
