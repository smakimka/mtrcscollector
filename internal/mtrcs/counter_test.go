package mtrcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric CounterMetric
		want   string
	}{
		{
			name:   "pozitive number",
			metric: CounterMetric{Name: "", Value: 1},
			want:   "1",
		},
		{
			name:   "negative number",
			metric: CounterMetric{Name: "", Value: -1},
			want:   "-1",
		},
		{
			name:   "negative zero",
			metric: CounterMetric{Name: "", Value: -0},
			want:   "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.metric.GetStringValue())
		})
	}
}
