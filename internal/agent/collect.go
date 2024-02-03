package agent

import (
	"context"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func CollectMetrics(ctx context.Context, s storage.Storage) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	UpdateMetrics(ctx, &m, s)
}

func UpdateMetrics(ctx context.Context, m *runtime.MemStats, s storage.Storage) {
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "Alloc", Value: float64(m.Alloc)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "BuckHashSys", Value: float64(m.BuckHashSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "Frees", Value: float64(m.Frees)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "GCCPUFraction", Value: m.GCCPUFraction})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "GCSys", Value: float64(m.GCSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapAlloc", Value: float64(m.HeapAlloc)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapIdle", Value: float64(m.HeapIdle)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapInuse", Value: float64(m.HeapInuse)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapObjects", Value: float64(m.HeapObjects)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapReleased", Value: float64(m.HeapReleased)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "HeapSys", Value: float64(m.HeapSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "LastGC", Value: float64(m.LastGC)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "Lookups", Value: float64(m.Lookups)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "MCacheInuse", Value: float64(m.MCacheInuse)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "MCacheSys", Value: float64(m.MCacheSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "MSpanInuse", Value: float64(m.MSpanInuse)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "MSpanSys", Value: float64(m.MSpanSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "Mallocs", Value: float64(m.Mallocs)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "NextGC", Value: float64(m.NextGC)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "NumForcedGC", Value: float64(m.NumForcedGC)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "NumGC", Value: float64(m.NumGC)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "OtherSys", Value: float64(m.OtherSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "StackInuse", Value: float64(m.StackInuse)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "StackSys", Value: float64(m.StackSys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "Sys", Value: float64(m.Sys)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "TotalAlloc", Value: float64(m.TotalAlloc)})

	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "RandomValue", Value: rand.Float64()})
	s.UpdateCounterMetric(ctx, model.CounterMetric{Name: "PollCount", Value: 1})
}

func CollectPSutilMetrics(ctx context.Context, s storage.Storage, errs chan<- error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		errs <- err
		return
	}

	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "TotalMemory", Value: float64(v.Total)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "FreeMemory", Value: float64(v.Free)})
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "CPUutilization1", Value: float64(v.UsedPercent)})
}
