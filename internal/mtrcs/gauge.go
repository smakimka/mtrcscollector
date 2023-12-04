package mtrcs

import "fmt"

type GaugeMetric struct {
	Name  string
	Value float64
}

func (m GaugeMetric) GetType() string {
	return "gauge"
}

func (m GaugeMetric) GetValue() float64 {
	return m.Value
}

func (m GaugeMetric) GetName() string {
	return m.Name
}

func (m GaugeMetric) GetStringValue() string {
	return fmt.Sprint(m.Value)
}
