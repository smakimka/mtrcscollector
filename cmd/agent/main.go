package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent"
	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/auth"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
	na           = "N/A"
)

func main() {
	if buildVersion == "" {
		buildVersion = na
	}
	if buildDate == "" {
		buildDate = na
	}
	if buildCommit == "" {
		buildCommit = na
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	cfg := config.NewConfig()
	logger.SetLevel(logger.Info)

	if cfg.Key != "" {
		auth.Init(cfg.Key)
	}

	s := storage.NewMemStorage()

	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s", cfg.Addr))

	if cfg.CryptoKeyPath != "" {
		if err := cfg.ReadCryptoKey(); err != nil {
			panic(err)
		}
	}

	if err := cfg.SetMyIP(); err != nil {
		panic(err)
	}

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
		go agent.Worker(ctx, *cfg, client, i+1, jobs, errs)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for range c {
			pollTicker.Stop()
			reportTicker.Stop()

			fmt.Println("Just a second, sending data...")
			agent.SendMetrics(context.Background(), cfg, s, jobs, errs)
			fmt.Println("Done!")
			os.Exit(0)
		}
	}()

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
