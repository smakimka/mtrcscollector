package mtrcs

const (
	Gauge   = iota
	Counter = iota
)

type Metric interface {
	GetType() int
}
