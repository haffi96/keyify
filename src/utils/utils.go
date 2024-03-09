package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
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
