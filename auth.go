package documentdb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Key struct {
	Key  string
	salt []byte
	err  error
}

func NewKey(key string) *Key {
	return &Key{Key: key}
}

var (
	enc = base64.StdEncoding
)

func (k *Key) Salt() ([]byte, error) {
	if len(k.salt) == 0 && k.err == nil {
		k.salt, k.err = enc.DecodeString(k.Key)
		if k.err != nil {
			if _, ok := k.err.(base64.CorruptInputError); ok {
				k.err = errors.New("base64 input is corrupt, check CosmosDB key.")
			}
		}
	}
	return k.salt, k.err
}

func authorize(str []byte, key *Key) (ret string, err error) {
	var (
		salt []byte
	)
	salt, err = key.Salt()

	if err != nil {
		return ret, err
	}

	hmac := hmac.New(sha256.New, salt)
	hmac.Write(str)
	b := hmac.Sum(nil)

	ret = enc.EncodeToString(b)
	return ret, nil
}
