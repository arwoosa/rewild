package helpers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func GetSignature(input string, key string) string {
	signingKey := []byte(key)
	h := hmac.New(sha1.New, signingKey)
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
