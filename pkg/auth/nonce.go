package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func NewNonce() (string, error) {
	nonceBytes := make([]byte, 32)

	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
