package mtrcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric GaugeMetric
		want   string
	}{
		{
			name:   "pozitive number",
			metric: GaugeMetric{Name: "", Value: 1.5},
			want:   "1.5",
		},
		{
			name:   "negative number",
			metric: GaugeMetric{Name: "", Value: -1.5},
			want:   "-1.5",
		},
		{
			name:   "negative zero",
			metric: GaugeMetric{Name: "", Value: -0},
			want:   "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.metric.GetStringValue())
		})
	}
}
