package main

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr                       string `env:"ADDRESS"`
	ReportIntervalVal          int    `env:"REPORT_INTERVAL"`
	PollIntervalVal            int    `env:"POLL_INTERVAL"`
	reportInterval             time.Duration
	pollInterval               time.Duration
	concurrentMetricsSendCount int
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
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Addr == "" {
		cfg.Addr = serverAddr
	}

	if cfg.PollIntervalVal == 0 {
		cfg.pollInterval = pollInteraval
	} else {
		cfg.pollInterval = time.Duration(cfg.PollIntervalVal) * time.Second
	}

	if cfg.ReportIntervalVal == 0 {
		cfg.reportInterval = reportInterval
	} else {
		cfg.reportInterval = time.Duration(cfg.ReportIntervalVal) * time.Second
	}

	cfg.concurrentMetricsSendCount = 10

	return cfg
}
