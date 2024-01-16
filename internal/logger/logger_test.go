package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	err := SetLevel(256)
	assert.Equal(t, err, ErrNoSuchLevel)
}
