package main

import (
	"fmt"
	"net/http"

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

	var s storage.Storage
	if cfg.StoreInterval == 0 {
		s = storage.NewSyncMemStorage(cfg.FileStoragePath)
	} else {
		s = storage.NewMemStorage()
	}

	if cfg.Restore {
		if err := s.Restore(cfg.FileStoragePath); err != nil {
			return err
		}
	}

	logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))
	return http.ListenAndServe(cfg.Addr, router.GetRouter(s))
}
