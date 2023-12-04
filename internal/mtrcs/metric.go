package mtrcs

const (
	Gauge   = iota
	Counter = iota
)

type Metric interface {
	GetType() string
	GetValue() float64
	GetName() string
	GetStringValue() string
}
