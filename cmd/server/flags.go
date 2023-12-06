package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr string `env:"ADDRESS"`
}

func parseFlags() *Config {
	var flagRunAddr string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "host:port to run on")
	flag.Parse()

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Addr == "" {
		cfg.Addr = flagRunAddr
	}

	return cfg
}
