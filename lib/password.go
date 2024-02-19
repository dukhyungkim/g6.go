package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

const (
	Pbkdf2CompatHashAlgorithm = "SHA256"
	Pbkdf2CompatIterations    = 12000
	Pbkdf2CompatSaltBytes     = 24
	Pbkdf2CompatHashBytes     = 24
)

func CreateHash(password string) string {
	salt := mustGenerateSalt(Pbkdf2CompatSaltBytes)
	algo := strings.ToLower(Pbkdf2CompatHashAlgorithm)
	iterations := Pbkdf2CompatIterations

	pbkdf2Bytes := pbkdf2.Key([]byte(password), []byte(salt), iterations, Pbkdf2CompatHashBytes, sha256.New)
	hash := base64.StdEncoding.EncodeToString(pbkdf2Bytes)

	return fmt.Sprintf("%s:%d:%s:%s", algo, iterations, salt, hash)
}

func mustGenerateSalt(bytes int) string {
	salt := make([]byte, bytes)
	_, err := rand.Read(salt)
	if err != nil {
		panic(fmt.Errorf("failed to generate salt: %w", err))
	}
	return base64.StdEncoding.EncodeToString(salt)
}
