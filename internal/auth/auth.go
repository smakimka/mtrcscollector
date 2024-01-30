package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

var (
	key []byte
)

func Enabled() bool {
	return key != nil
}

func Init(newKey string) {
	key = []byte(newKey)
}

func Sign(data []byte) []byte {
	haser := GetHasher()
	return haser.Sum(data)
}

func Check(originalSign []byte, data []byte) (bool, error) {
	hasher := GetHasher()
	sign := hasher.Sum(data)
	return hmac.Equal(originalSign, sign), nil
}

func GetHasher() hash.Hash {
	return hmac.New(sha256.New, key)
}
