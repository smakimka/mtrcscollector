package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig(t *testing.T) {
	defaultValues := struct {
		addr            string
		fileStoragePath string
		storeInterval   int
		restore         bool
	}{
		"localhost:8080",
		"/tmp/metrics-db.json",
		300,
		true,
	}

	cfg := NewConfig()
	assert.Equal(t, defaultValues.addr, cfg.Addr)
	assert.Equal(t, defaultValues.storeInterval, cfg.StoreInterval)
	assert.Equal(t, defaultValues.fileStoragePath, cfg.FileStoragePath)
	assert.Equal(t, defaultValues.restore, cfg.Restore)
}
