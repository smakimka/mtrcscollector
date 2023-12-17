package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent"
	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	cfg := config.NewConfig()

	logger := log.New(os.Stdout, "", 5)

	s := storage.NewMemStorage().
		WithLogger(logger)

	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s", cfg.Addr))

	// инициализация метрик
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	agent.UpdateMetrics(&m, s, logger)
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "LastPollCount", Value: 0})

	run(cfg, s, logger, client)
}

func run(cfg *config.Config, s storage.Storage, l *log.Logger, client *resty.Client) {
	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	errChan := make(chan error)

	for {
		select {
		case <-pollTicker.C:
			go agent.CollectMetrics(cfg, s, l)
		case <-reportTicker.C:
			go agent.SendMetrics(cfg, s, l, client, errChan)
		case err := <-errChan:
			panic(err)
		}
	}
}
