package utils

import (
	"crypto/rand"
	"encoding/hex"
)

const length = 8

func GenerateRequestID() (string, error) {
	requestID := make([]byte, length)
	_, err := rand.Read(requestID)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(requestID), nil
}
