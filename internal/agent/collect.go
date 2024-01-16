package agent

import (
	"math/rand"
	"runtime"

	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func CollectMetrics(cfg *config.Config, s storage.Storage) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	UpdateMetrics(&m, s)
}

func UpdateMetrics(m *runtime.MemStats, s storage.Storage) {
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "Alloc", Value: float64(m.Alloc)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "BuckHashSys", Value: float64(m.BuckHashSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "Frees", Value: float64(m.Frees)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "GCCPUFraction", Value: m.GCCPUFraction})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "GCSys", Value: float64(m.GCSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapAlloc", Value: float64(m.HeapAlloc)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapIdle", Value: float64(m.HeapIdle)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapInuse", Value: float64(m.HeapInuse)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapObjects", Value: float64(m.HeapObjects)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapReleased", Value: float64(m.HeapReleased)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "HeapSys", Value: float64(m.HeapSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "LastGC", Value: float64(m.LastGC)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "Lookups", Value: float64(m.Lookups)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "MCacheInuse", Value: float64(m.MCacheInuse)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "MCacheSys", Value: float64(m.MCacheSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "MSpanInuse", Value: float64(m.MSpanInuse)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "MSpanSys", Value: float64(m.MSpanSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "Mallocs", Value: float64(m.Mallocs)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "NextGC", Value: float64(m.NextGC)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "NumForcedGC", Value: float64(m.NumForcedGC)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "NumGC", Value: float64(m.NumGC)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "OtherSys", Value: float64(m.OtherSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "StackInuse", Value: float64(m.StackInuse)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "StackSys", Value: float64(m.StackSys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "Sys", Value: float64(m.Sys)})
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "TotalAlloc", Value: float64(m.TotalAlloc)})

	s.UpdateGaugeMetric(model.GaugeMetric{Name: "RandomValue", Value: rand.Float64()})
	s.UpdateCounterMetric(model.CounterMetric{Name: "PollCount", Value: 1})
}
