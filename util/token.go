package util

import (
	"crypto/rand"
	"encoding/hex"
)

func TokenHex(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}
