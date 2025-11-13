package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecureBytes generates a slice of n cryptographically secure random bytes
func GenerateSecureBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("length must be greater than 0")
	}
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read from crypto/rand: %v", err)
	}
	return b, nil
}

// GenerateSecureString generates a base64 encoded string of n cryptographically secure random bytes
func GenerateSecureString(n int) (string, error) {
	bytes, err := GenerateSecureBytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
