package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		metric GaugeMetric
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

func TestGetGaugeValue(t *testing.T) {
	tests := []struct {
		name   string
		metric GaugeMetric
		want   float64
	}{
		{
			name:   "data type is float64",
			metric: GaugeMetric{Name: "test", Value: 1.5},
			want:   1.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.metric.GetValue())
		})
	}
}
