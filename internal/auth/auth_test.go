package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignCheck(t *testing.T) {
	data := []byte("test string")
	sign := Sign(data)

	res, err := Check(sign, data)
	assert.NoError(t, err)
	assert.True(t, res)
}
