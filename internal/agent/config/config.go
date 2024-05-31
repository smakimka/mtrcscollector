package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"net"
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
	MyIP           string
	GRPC           bool
}

type JsonConfig struct {
	Addr           string `json:"addr"`
	CryptoKey      string `json:"crypto_key"`
	Key            string `json:"key"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	RateLimit      int    `json:"rate_limit"`
	GRPC           string `json:"grpc"`
}

type EnvParams struct {
	Config         string `env:"CONFIG"`
	Addr           string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	CryptoKeyPath  string `env:"CRYPTO_KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	GRPC           string `env:"GRPC"`
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

func (c *Config) SetMyIP() error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				c.MyIP = ipnet.IP.String()
				return nil
			}
		}
	}
	return nil
}

func parseFlags() *Config {
	var serverAddr string
	var flagReportInterval int
	var flagPollInteraval int
	var rateLimit int
	var flagGRPC bool

	var flagKey string
	var flagCryptoKey string

	var flagConfig string

	flag.StringVar(&flagConfig, "c", "{}", "config in json format")

	flag.StringVar(&serverAddr, "a", "localhost:8080", "server addres without http://")

	flag.IntVar(&flagReportInterval, "r", 10, "metrics sending period (in seconds)")
	flag.IntVar(&flagPollInteraval, "p", 2, "metrics updqating period (in seconds)")
	reportInterval := time.Duration(flagReportInterval) * time.Second
	pollInteraval := time.Duration(flagPollInteraval) * time.Second
	flag.StringVar(&flagKey, "k", "", "auth key string")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "path to public key file")

	flag.IntVar(&rateLimit, "l", 1, "number of max concurrent request")

	flag.BoolVar(&flagGRPC, "g", false, "grpc or not")
	flag.Parse()

	var jsonCfg JsonConfig
	err := json.Unmarshal([]byte(flagConfig), &jsonCfg)
	if err != nil {
		panic(err)
	}

	cfg := &Config{}
	envParams := &EnvParams{}
	err = env.Parse(envParams)
	if err != nil {
		panic(err)
	}

	readJson(cfg, &jsonCfg)

	if envParams.Addr == "" {
		if serverAddr != "localhost:8080" {
			cfg.Addr = serverAddr
		} else {
			if cfg.Addr == "" {
				cfg.Addr = serverAddr
			}
		}
	} else {
		cfg.Addr = envParams.Addr
	}

	if envParams.PollInterval == 0 {
		if flagPollInteraval != 2 {
			cfg.PollInterval = pollInteraval
		} else {
			if cfg.PollInterval == 0 {
				cfg.PollInterval = pollInteraval
			}
		}
	} else {
		cfg.PollInterval = time.Duration(envParams.PollInterval) * time.Second
	}

	if envParams.ReportInterval == 0 {
		if flagReportInterval != 10 {
			cfg.ReportInterval = reportInterval
		} else {
			if cfg.ReportInterval == 0 {
				cfg.ReportInterval = reportInterval
			}
		}
	} else {
		cfg.ReportInterval = time.Duration(envParams.ReportInterval) * time.Second
	}

	if envParams.RateLimit == 0 {
		if rateLimit != 1 {
			cfg.RateLimit = rateLimit
		} else {
			if cfg.RateLimit == 0 {
				cfg.RateLimit = rateLimit
			}
		}
	} else {
		cfg.RateLimit = envParams.RateLimit
	}

	if envParams.Key == "" {
		if flagKey != "" {
			cfg.Key = flagKey
		} else {
			if cfg.Key == "" {
				cfg.Key = flagKey
			}
		}
	} else {
		cfg.Key = envParams.Key
	}

	if envParams.CryptoKeyPath == "" {
		if flagCryptoKey != "" {
			cfg.CryptoKeyPath = flagCryptoKey
		} else {
			if cfg.CryptoKeyPath == "" {
				cfg.CryptoKeyPath = flagCryptoKey
			}
		}
		cfg.CryptoKeyPath = flagCryptoKey
	}

	if envParams.GRPC == "" {
		if !flagGRPC {
			cfg.GRPC = flagGRPC
		} else {
			if !cfg.GRPC {
				cfg.GRPC = flagGRPC
			}
		}
		cfg.GRPC = flagGRPC
	}

	return cfg
}

func readJson(cfg *Config, jsonCfg *JsonConfig) {
	if jsonCfg.Addr != "" {
		cfg.Addr = jsonCfg.Addr
	}
	if jsonCfg.CryptoKey != "" {
		cfg.CryptoKeyPath = jsonCfg.CryptoKey
	}
	if jsonCfg.Key != "" {
		cfg.Key = jsonCfg.Key
	}
	if jsonCfg.RateLimit != 0 {
		cfg.RateLimit = jsonCfg.RateLimit
	}
	if jsonCfg.PollInterval != 0 {
		cfg.PollInterval = time.Duration(jsonCfg.PollInterval) * time.Second
	}
	if jsonCfg.ReportInterval != 0 {
		cfg.ReportInterval = time.Duration(jsonCfg.ReportInterval) * time.Second
	}
}
