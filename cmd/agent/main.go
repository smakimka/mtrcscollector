package main

import (
	"context"
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

	ctx := context.Background()
	// инициализация метрик
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	agent.UpdateMetrics(ctx, &m, s)
	s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "LastPollCount", Value: 0})

	run(ctx, cfg, s, client)
}

func run(ctx context.Context, cfg *config.Config, s storage.Storage, client *resty.Client) {
	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	jobs := make(chan model.MetricsData, cfg.RateLimit)
	errs := make(chan error)

	for i := 0; i < cfg.RateLimit; i++ {
		go agent.Worker(ctx, client, i+1, jobs, errs)
	}

	for {
		select {
		case <-pollTicker.C:
			go agent.CollectMetrics(ctx, s)
			go agent.CollectPSutilMetrics(ctx, s, errs)
		case <-reportTicker.C:
			go agent.SendMetrics(ctx, cfg, s, jobs, errs)
		case err := <-errs:
			fmt.Println(err)
		}
	}
}
