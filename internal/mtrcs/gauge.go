package mtrcs

type GaugeMetric struct {
	Name  string
	Value float64
}

func (m GaugeMetric) GetType() int {
	return Gauge
}
