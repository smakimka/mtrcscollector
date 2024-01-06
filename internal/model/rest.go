package model

import (
	"errors"
	"net/http"
)

var ErrMissingFields = errors.New("missing some of required fields")
var ErrWrongMetricKind = errors.New("wrong metric kind")

type MetricsData struct {
	Name  string   `json:"id"`              // имя метрики
	Kind  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *MetricsData) Bind(r *http.Request) error {
	if m.Name == "" || m.Kind == "" {
		return ErrMissingFields
	}

	if m.Kind != Gauge && m.Kind != Counter {
		return ErrWrongMetricKind
	}

	return nil
}

type Response struct {
	Ok     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}
