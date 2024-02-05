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
	RateLimit      int
}

type EnvParams struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig() *Config {
	return parseFlags()
}

func parseFlags() *Config {
	var serverAddr string
	var flagReportInterval int
	var flagPollInteraval int
	var rateLimit int

	flag.StringVar(&serverAddr, "a", "localhost:8080", "server addres without http://")

	flag.IntVar(&flagReportInterval, "r", 10, "metrics sending period (in seconds)")
	flag.IntVar(&flagPollInteraval, "p", 2, "metrics updqating period (in seconds)")
	reportInterval := time.Duration(flagReportInterval) * time.Second
	pollInteraval := time.Duration(flagPollInteraval) * time.Second

	flag.IntVar(&rateLimit, "l", 1, "number of max concurrent request")

	flag.Parse()

	cfg := &Config{}
	envParams := &EnvParams{}
	err := env.Parse(envParams)
	if err != nil {
		panic(err)
	}

	if envParams.Addr == "" {
		cfg.Addr = serverAddr
	} else {
		cfg.Addr = envParams.Addr
	}

	if envParams.PollInterval == 0 {
		cfg.PollInterval = pollInteraval
	} else {
		cfg.PollInterval = time.Duration(envParams.PollInterval) * time.Second
	}

	if envParams.ReportInterval == 0 {
		cfg.ReportInterval = reportInterval
	} else {
		cfg.ReportInterval = time.Duration(envParams.ReportInterval) * time.Second
	}

	if envParams.RateLimit == 0 {
		cfg.RateLimit = rateLimit
	} else {
		cfg.RateLimit = envParams.RateLimit
	}

	return cfg
}
