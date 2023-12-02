package mtrcs

type CounterMetric struct {
	Name  string
	Value int64
}

func (m CounterMetric) GetType() int {
	return Gauge
}
