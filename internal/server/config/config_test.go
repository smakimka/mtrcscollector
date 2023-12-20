package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig(t *testing.T) {
	defaultValues := struct {
		addr string
	}{
		"localhost:8080",
	}

	cfg := NewConfig()
	assert.Equal(t, defaultValues.addr, cfg.Addr)
}
