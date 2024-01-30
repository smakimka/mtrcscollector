package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smakimka/mtrcscollector/internal/auth"
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

	var s storage.Storage
	if cfg.DatabaseDSN == "" {
		storage, err := initSyncStorage(cfg)
		if err != nil {
			return err
		}
		s = storage
	} else {
		pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return err
		}
		defer pool.Close()

		s, err = storage.NewPGStorage(ctx, pool)
		if err != nil {
			return err
		}
	}

	if cfg.Key != "" {
		auth.Init(cfg.Key)
	}

	logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))
	return http.ListenAndServe(cfg.Addr, router.GetRouter(s))
}

func initSyncStorage(cfg *config.Config) (storage.SyncStorage, error) {
	var s storage.SyncStorage
	if cfg.StoreInterval == 0 {
		s = storage.NewSyncMemStorage(cfg.FileStoragePath)
	} else {
		s = storage.NewMemStorage()
		go saveMetrics(s, cfg)
	}

	if cfg.Restore {
		if err := s.Restore(cfg.FileStoragePath); err != nil {
			return nil, err
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

	return s, nil
}

func saveMetrics(s storage.SyncStorage, cfg *config.Config) {
	saveTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

	for range saveTicker.C {
		go s.Save(cfg.FileStoragePath)
	}
}
