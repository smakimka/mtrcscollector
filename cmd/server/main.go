package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smakimka/mtrcscollector/internal/auth"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/server/config"
	"github.com/smakimka/mtrcscollector/internal/server/grpc"
	"github.com/smakimka/mtrcscollector/internal/server/router"
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

	if cfg.CryptoKeyPath != "" {
		if err := cfg.ReadCryptoKey(); err != nil {
			panic(err)
		}
	}

	if cfg.TrustedSubnetString != "" {
		if err := cfg.ParseCIDR(); err != nil {
			panic(err)
		}
	}

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

	if cfg.StartAsGRPC {
		listen, err := net.Listen("tcp", cfg.Addr)
		if err != nil {
			return err
		}

		logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))
		server := grpc.NewServer(cfg, s)
		if err := server.Serve(listen); err != nil {
			return err
		}
	}

	logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))
	return http.ListenAndServe(cfg.Addr, router.GetRouter(s, cfg.CryptoKey, cfg.TrustedSubnet))
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
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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
