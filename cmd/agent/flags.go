package main

import (
	"flag"
	"time"
)

var (
	serverAddr                 string
	flagReportInterval         int
	flagPollInteraval          int
	reportInterval             time.Duration
	pollInteraval              time.Duration
	concurrentMetricsSendCount = 10
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server addres without http://")

	flag.IntVar(&flagReportInterval, "r", 10, "metrics sending period (in seconds)")
	flag.IntVar(&flagPollInteraval, "p", 2, "metrics updqating period (in seconds)")
	reportInterval = time.Duration(flagReportInterval) * time.Second
	pollInteraval = time.Duration(flagPollInteraval) * time.Second

	flag.Parse()
}
