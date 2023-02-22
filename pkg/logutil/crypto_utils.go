package logutil

import (
	"encoding/base64"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func CryptoKeyToBytes(key crypto.Key) []byte {
	keyBytes, err := key.Raw()
	if err != nil {
		return []byte{0}
	}

	return keyBytes
}

func CryptoKeyToBase64(key crypto.Key) string {
	bytes := CryptoKeyToBytes(key)
	return base64.StdEncoding.EncodeToString(bytes)
}
