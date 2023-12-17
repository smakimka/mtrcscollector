package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr           string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

type EnvParams struct {
	addr           string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func NewConfig() *Config {
	return parseFlags()
}

func parseFlags() *Config {
	var serverAddr string
	var flagReportInterval int
	var flagPollInteraval int

	flag.StringVar(&serverAddr, "a", "localhost:8080", "server addres without http://")

	flag.IntVar(&flagReportInterval, "r", 10, "metrics sending period (in seconds)")
	flag.IntVar(&flagPollInteraval, "p", 2, "metrics updqating period (in seconds)")
	reportInterval := time.Duration(flagReportInterval) * time.Second
	pollInteraval := time.Duration(flagPollInteraval) * time.Second

	flag.Parse()

	cfg := &Config{}
	envParams := &EnvParams{}
	err := env.Parse(envParams)
	if err != nil {
		panic(err)
	}

	if envParams.addr == "" {
		cfg.Addr = serverAddr
	}

	if envParams.pollInterval == 0 {
		cfg.PollInterval = pollInteraval
	} else {
		cfg.PollInterval = time.Duration(envParams.pollInterval) * time.Second
	}

	if envParams.reportInterval == 0 {
		cfg.ReportInterval = reportInterval
	} else {
		cfg.ReportInterval = time.Duration(envParams.reportInterval) * time.Second
	}

	return cfg
}
