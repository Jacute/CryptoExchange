package utils

import (
	"encoding/hex"
	"math/rand"
)

func GenerateToken(length int) string {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = byte(rand.Intn(256))
	}

	return hex.EncodeToString(bytes)
}
