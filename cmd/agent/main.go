package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	parseFlags()

	logger := log.New(os.Stdout, "", 5)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	// инициализация метрик
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	updateMetrics(&m, s, logger)
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "LastPollCount", Value: 0})

	var wg sync.WaitGroup

	wg.Add(1)
	go collectMetrics(&wg, s, logger)
	go sendMetrics(&wg, s, client, logger)
	wg.Wait()
}
