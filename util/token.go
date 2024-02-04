package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func TokenHex(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even")
	}

	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}
