package mtrcs

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric interface {
	GetType() string
	GetValue() float64
	GetName() string
	GetStringValue() string
}
