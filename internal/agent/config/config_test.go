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
	assert.Equal(t, cfg.Addr, defaultValues.addr)
	assert.Equal(t, cfg.ReportInterval, defaultValues.reportInterval)
	assert.Equal(t, cfg.PollInterval, defaultValues.pollInterval)
}
