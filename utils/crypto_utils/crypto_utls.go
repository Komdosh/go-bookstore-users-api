package crypto_utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetSha1(input string) string {
	hash := sha1.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}
