package documentdb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func authorize(str, key string) (string, error) {
	enc := base64.StdEncoding
	salt, err := enc.DecodeString(key)
	if err != nil {
		return "", err
	}
	hmac := hmac.New(sha256.New, salt)
	hmac.Write([]byte(str))
	b := hmac.Sum(nil)
	return enc.EncodeToString(b), nil
}
