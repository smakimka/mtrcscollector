package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	cfg := parseFlags()

	logger := log.New(os.Stdout, "", 5)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s", cfg.Addr))

	// инициализация метрик
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	updateMetrics(&m, s, logger)
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "LastPollCount", Value: 0})

	var wg sync.WaitGroup

	wg.Add(1)
	go collectMetrics(cfg, &wg, s, logger)
	go sendMetrics(cfg, &wg, s, client, logger)
	wg.Wait()
}
