package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr            string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	CryptoKeyPath   string `env:"CRYPTO_KEY"`
	CryptoKey       *rsa.PrivateKey
	StoreInterval   int  `env:"STORE_INTERVAL"`
	Restore         bool `env:"RESTORE"`
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

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "host:port to run on")
	flag.IntVar(&flagStoreInterval, "i", 300, "state save interval (in seconds)")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "temp file to save state to (if emtpy no saves are done)")
	flag.BoolVar(&flagRestore, "r", true, "load with saved data or not")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database dsn string")
	flag.StringVar(&flagKey, "k", "", "auth key string")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "path to a private key file")

	flag.Parse()

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Addr == "" {
		cfg.Addr = flagRunAddr
	}

	if os.Getenv("STORE_INTERVAL") == "" {
		cfg.StoreInterval = flagStoreInterval
	}

	if os.Getenv("FILE_STORAGE_PATH") == "" {
		cfg.FileStoragePath = flagStoragePath
	}

	if os.Getenv("RESTORE") == "" {
		cfg.Restore = flagRestore
	}

	if os.Getenv("DATABASE_DSN") == "" {
		cfg.DatabaseDSN = flagDatabaseDSN
	}

	if os.Getenv("KEY") == "" {
		cfg.Key = flagKey
	}

	if os.Getenv("CryptoKey") == "" {
		cfg.CryptoKeyPath = flagCryptoKey
	}

	return cfg
}
