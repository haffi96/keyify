package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	mathRand "math/rand"
	"time"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HashString(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func GenerateRandomId(prefix string) string {
	// Length of the random string (excluding prefix)
	randomStringLength := 16

	// Characters to use for the random string
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate random string directly
	b := make([]byte, randomStringLength)
	for i := range b {
		b[i] = chars[mathRand.Intn(len(chars))] // Select a random character
	}
	randomString := string(b)

	// Combine the prefix and random string
	keyId := prefix + randomString

	return keyId
}

func GenerateApiKey(preifx string) (string, error) {
	// Generate random API key
	apiKeyBytes, err := GenerateRandomBytes(36)
	if err != nil {
		return "", err
	}
	apiKeyString := base64.URLEncoding.EncodeToString(apiKeyBytes)
	return fmt.Sprintf("%s%s", preifx, apiKeyString), nil
}

// TimeNow returns the current time formatted as a string in RFC3339 format.
func TimeNow() string {
	return time.Now().Format(time.RFC3339)
}
