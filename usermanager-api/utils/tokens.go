package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateVerificationToken generates a random verification token
func GenerateVerificationToken() (string, error) {
	// Create a byte array to hold the random bytes
	b := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode the random bytes to a base64 string
	return base64.URLEncoding.EncodeToString(b), nil
}
