package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr            string `env:"ADDRESS" json:"addr"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	Key             string `env:"KEY" json:"key"`
	CryptoKeyPath   string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKey       *rsa.PrivateKey
	StoreInterval   int  `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool `env:"RESTORE" json:"restore"`
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

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	c.CryptoKey = key

	return nil
}

func parseFlags() *Config {
	var flagRunAddr string
	var flagStoreInterval int
	var flagStoragePath string
	var flagRestore bool
	var flagDatabaseDSN string
	var flagKey string
	var flagCryptoKey string
	var flagJsonConfig string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "host:port to run on")
	flag.IntVar(&flagStoreInterval, "i", 300, "state save interval (in seconds)")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "temp file to save state to (if emtpy no saves are done)")
	flag.BoolVar(&flagRestore, "r", true, "load with saved data or not")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database dsn string")
	flag.StringVar(&flagKey, "k", "", "auth key string")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "path to a private key file")
	flag.StringVar(&flagJsonConfig, "c", "{}", "config in json format")

	flag.Parse()

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	var jsonCfg Config
	err = json.Unmarshal([]byte(flagJsonConfig), &jsonCfg)
	if err != nil {
		panic(err)
	}

	if cfg.Addr == "" {
		if flagRunAddr != "localhost:8080" {
			cfg.Addr = flagRunAddr
		} else {
			if jsonCfg.Addr != "" {
				cfg.Addr = jsonCfg.Addr
			} else {
				cfg.Addr = flagRunAddr
			}
		}
	}

	if os.Getenv("STORE_INTERVAL") == "" {
		if flagStoreInterval != 300 {
			cfg.StoreInterval = flagStoreInterval
		} else {
			if jsonCfg.StoreInterval != 0 {
				cfg.StoreInterval = jsonCfg.StoreInterval
			} else {
				cfg.StoreInterval = flagStoreInterval
			}
		}
	}

	if os.Getenv("FILE_STORAGE_PATH") == "" {
		if flagStoragePath != "/tmp/metrics-db.json" {
			cfg.FileStoragePath = flagStoragePath
		} else {
			if jsonCfg.FileStoragePath != "" {
				cfg.FileStoragePath = jsonCfg.FileStoragePath
			} else {
				cfg.FileStoragePath = flagStoragePath
			}
		}
	}

	if os.Getenv("RESTORE") == "" {
		if !flagRestore {
			cfg.Restore = flagRestore
		} else {
			if jsonCfg.Restore {
				cfg.Restore = jsonCfg.Restore
			} else {
				cfg.Restore = flagRestore
			}
		}
	}

	if os.Getenv("DATABASE_DSN") == "" {
		if flagDatabaseDSN != "" {
			cfg.DatabaseDSN = flagDatabaseDSN
		} else {
			if jsonCfg.DatabaseDSN != "" {
				cfg.DatabaseDSN = jsonCfg.DatabaseDSN
			} else {
				cfg.DatabaseDSN = flagDatabaseDSN
			}
		}
	}

	if os.Getenv("KEY") == "" {
		if flagKey != "" {
			cfg.Key = flagKey
		} else {
			if jsonCfg.Key != "" {
				cfg.Key = jsonCfg.Key
			} else {
				cfg.Key = flagKey
			}
		}
	}

	if os.Getenv("CryptoKey") == "" {
		if flagCryptoKey != "" {
			cfg.CryptoKeyPath = flagCryptoKey
		} else {
			if jsonCfg.CryptoKeyPath != "" {
				cfg.CryptoKeyPath = jsonCfg.CryptoKeyPath
			} else {
				cfg.CryptoKeyPath = flagCryptoKey
			}
		}
	}

	return cfg
}
