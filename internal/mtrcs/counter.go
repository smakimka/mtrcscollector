package mtrcs

import "fmt"

type CounterMetric struct {
	Name  string
	Value int64
}

func (m CounterMetric) GetType() string {
	return "counter"
}

func (m CounterMetric) GetValue() float64 {
	return float64(m.Value)
}

func (m CounterMetric) GetName() string {
	return m.Name
}

func (m CounterMetric) GetStringValue() string {
	return fmt.Sprint(m.Value)
}
