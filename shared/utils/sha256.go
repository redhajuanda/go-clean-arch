package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))

	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}
