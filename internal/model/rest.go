package model

import (
	"errors"
	"net/http"
)

var ErrMissingFields = errors.New("missing some of required fields")
var ErrWrongMetricKind = errors.New("wrong metric kind")

type MetricData struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Name  string   `json:"id"`
	Kind  string   `json:"type"`
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
	Detail string `json:"detail,omitempty"`
	Ok     bool   `json:"ok"`
}
