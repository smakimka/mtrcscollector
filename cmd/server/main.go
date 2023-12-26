package main

import (
	"fmt"
	"net/http"

	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/server/config"
	"github.com/smakimka/mtrcscollector/internal/server/router"
)

func main() {
	cfg := config.NewConfig()

	if err := run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config) error {
	logger.SetLevel(logger.Info)
	logger.Log.Info().Msg(fmt.Sprintf("Running server on %s", cfg.Addr))

	return http.ListenAndServe(cfg.Addr, router.GetRouter())
}
