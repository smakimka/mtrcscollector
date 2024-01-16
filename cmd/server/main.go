package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/server/config"
	"github.com/smakimka/mtrcscollector/internal/server/router"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	cfg := config.NewConfig()
	if err := run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config) error {
	logger.SetLevel(logger.Info)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer pool.Close()

	var s storage.Storage
	if cfg.StoreInterval == 0 {
		s = storage.NewSyncMemStorage(cfg.FileStoragePath)
	} else {
		s = storage.NewMemStorage()
		go saveMetrics(s, cfg)
	}

	if cfg.Restore {
		if err := s.Restore(cfg.FileStoragePath); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("Just a second, saving data...")
			s.Save(cfg.FileStoragePath)
			fmt.Println("Done!")
			os.Exit(0)
		}
	}()

	logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))
	return http.ListenAndServe(cfg.Addr, router.GetRouter(s, pool))
}

func saveMetrics(s storage.Storage, cfg *config.Config) {
	saveTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

	for range saveTicker.C {
		go s.Save(cfg.FileStoragePath)
	}
}
