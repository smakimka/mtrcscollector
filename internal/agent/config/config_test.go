package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig(t *testing.T) {
	defaultValues := struct {
		addr           string
		pollInterval   time.Duration
		reportInterval time.Duration
	}{
		"localhost:8080",
		2 * time.Second,
		10 * time.Second,
	}

	cfg := NewConfig()
	assert.Equal(t, defaultValues.addr, cfg.Addr)
	assert.Equal(t, defaultValues.reportInterval, cfg.ReportInterval)
	assert.Equal(t, defaultValues.pollInterval, cfg.PollInterval)
}
