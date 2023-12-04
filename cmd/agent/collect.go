package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func collectMetrics(wg *sync.WaitGroup, s storage.Storage, logger *log.Logger) {
	defer wg.Done()

	for {
		m := runtime.MemStats{}
		runtime.ReadMemStats(&m)

		updateMetrics(&m, s, logger)

		time.Sleep(pollInteraval)
	}
}

func updateMetrics(m *runtime.MemStats, s storage.Storage, logger *log.Logger) {
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "Alloc", Value: float64(m.Alloc)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "BuckHashSys", Value: float64(m.BuckHashSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "GCCPUFraction", Value: m.GCCPUFraction})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "GCSys", Value: float64(m.GCSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapAlloc", Value: float64(m.HeapAlloc)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapIdle", Value: float64(m.HeapIdle)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapInuse", Value: float64(m.HeapInuse)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapObjects", Value: float64(m.HeapObjects)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapReleased", Value: float64(m.HeapReleased)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "HeapSys", Value: float64(m.HeapSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "LastGC", Value: float64(m.LastGC)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "Lookups", Value: float64(m.Lookups)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "MCacheInuse", Value: float64(m.MCacheInuse)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "MCacheSys", Value: float64(m.MCacheSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "MSpanInuse", Value: float64(m.MSpanInuse)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "MSpanSys", Value: float64(m.MSpanSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "Mallocs", Value: float64(m.Mallocs)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "NextGC", Value: float64(m.NextGC)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "NumForcedGC", Value: float64(m.NumForcedGC)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "NumGC", Value: float64(m.NumGC)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "OtherSys", Value: float64(m.OtherSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "StackInuse", Value: float64(m.StackInuse)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "StackSys", Value: float64(m.StackSys)})
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "TotalAlloc", Value: float64(m.TotalAlloc)})

	s.UpdateMetric(mtrcs.GaugeMetric{Name: "RandomValue", Value: rand.Float64()})
	s.UpdateMetric(mtrcs.CounterMetric{Name: "PollCount", Value: 1})
}
