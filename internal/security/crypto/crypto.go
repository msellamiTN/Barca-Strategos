package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomKey generates a cryptographically random key.
func GenerateRandomKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
