package model

import (
	"errors"
	"net/http"
)

var ErrMissingFields = errors.New("missing some of required fields")
var ErrWrongMetricKind = errors.New("wrong metric kind")

type MetricData struct {
	Name  string   `json:"id"`              // имя метрики
	Kind  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *MetricData) Bind(r *http.Request) error {
	if m.Name == "" || m.Kind == "" {
		return ErrMissingFields
	}

	if m.Kind != Gauge && m.Kind != Counter {
		return ErrWrongMetricKind
	}

	return nil
}

type MetricsData []MetricData

func (d MetricsData) Bind(r *http.Request) error {
	for _, metricData := range d {
		err := metricData.Bind(r)
		if err != nil {
			return err
		}
	}
	return nil
}

type Response struct {
	Ok     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}
