package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent"
	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/auth"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	cfg := config.NewConfig()
	logger.SetLevel(logger.Info)

	if cfg.Key != "" {
		auth.Init(cfg.Key)
	}

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

	errChan := make(chan error)

	for {
		select {
		case <-pollTicker.C:
			go agent.CollectMetrics(ctx, cfg, s)
		case <-reportTicker.C:
			go agent.SendMetrics(ctx, cfg, s, client, errChan)
		case err := <-errChan:
			fmt.Println(err)
		}
	}
}
