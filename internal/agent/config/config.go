package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr           string
	CryptoKeyPath  string
	CryptoKey      *rsa.PublicKey
	Key            string
	ReportInterval time.Duration
	PollInterval   time.Duration
	RateLimit      int
}

type EnvParams struct {
	Addr           string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	CryptoKeyPath  string `env:"CRYPTO_KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig() *Config {
	return parseFlags()
}

var ErrNokey = errors.New("key file doesn't contain key'")

func (c *Config) ReadCryptoKey() error {
	data, err := os.ReadFile(c.CryptoKeyPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return ErrNokey
	}

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}

	c.CryptoKey = key

	return nil
}

func parseFlags() *Config {
	var serverAddr string
	var flagReportInterval int
	var flagPollInteraval int
	var rateLimit int
	var flagKey string
	var flagCryptoKey string

	flag.StringVar(&serverAddr, "a", "localhost:8080", "server addres without http://")

	flag.IntVar(&flagReportInterval, "r", 10, "metrics sending period (in seconds)")
	flag.IntVar(&flagPollInteraval, "p", 2, "metrics updqating period (in seconds)")
	reportInterval := time.Duration(flagReportInterval) * time.Second
	pollInteraval := time.Duration(flagPollInteraval) * time.Second
	flag.StringVar(&flagKey, "k", "", "auth key string")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "path to public key file")

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

	if envParams.Key == "" {
		cfg.Key = flagKey
	} else {
		cfg.Key = envParams.Key
	}

	if envParams.CryptoKeyPath == "" {
		cfg.CryptoKey = flagCryptoKey
	}

	return cfg
}
