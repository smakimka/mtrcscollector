package storage_test

import (
	"context"
	"fmt"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func Example() {
	ctx := context.Background()
	s := storage.NewMemStorage()

	counterDelta := int64(10)
	gaugeValue := float64(10)
	s.UpdateMetrics(ctx, model.MetricsData{
		model.MetricData{
			Name:  "testCounter",
			Kind:  model.Counter,
			Delta: &counterDelta,
		},
		model.MetricData{
			Name:  "testGauge",
			Kind:  model.Gauge,
			Value: &gaugeValue,
		},
	})

	counterMetrics, err := s.GetAllCounterMetrics(ctx)
	fmt.Println(counterMetrics, err)

	gaugeMetrics, err := s.GetAllCounterMetrics(ctx)
	fmt.Println(gaugeMetrics, err)

	counterMetric, err := s.GetCounterMetric(ctx, "testCounter")
	fmt.Println(counterMetric, err)

	gaugeMetric, err := s.GetGaugeMetric(ctx, "testGauge")
	fmt.Println(gaugeMetric, err)
}
