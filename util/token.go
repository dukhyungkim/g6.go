package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func TokenHex(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}

func TokenURLSafe(length int) (string, error) {
	tok, err := TokenHex(length)
	if err != nil {
		return "", err
	}

	encodedTok := base64.URLEncoding.EncodeToString([]byte(tok))
	return encodedTok, nil
}
