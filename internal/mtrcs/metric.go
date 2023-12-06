package mtrcs

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric interface {
	GetType() string
	GetValue() interface{}
	GetName() string
	GetStringValue() string
}
