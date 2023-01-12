package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func PayloadSignature(timestamp, key string) string {
	mac := hmac.New(sha256.New, []byte(key))

	c := fmt.Sprintf("%s\n%s", timestamp, key)
	mac.Write([]byte(c))

	h := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(h)
}
