package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent"
	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	cfg := config.NewConfig()
	logger.SetLevel(logger.Info)

	s := storage.NewMemStorage()

	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s", cfg.Addr))

	// инициализация метрик
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	agent.UpdateMetrics(&m, s)
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "LastPollCount", Value: 0})

	run(cfg, s, client)
}

func run(cfg *config.Config, s storage.Storage, client *resty.Client) {
	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	errChan := make(chan error)

	for {
		select {
		case <-pollTicker.C:
			go agent.CollectMetrics(cfg, s)
		case <-reportTicker.C:
			go agent.SendMetrics(cfg, s, client, errChan)
		case err := <-errChan:
			fmt.Println(err)
		}
	}
}
