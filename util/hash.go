package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// salt를 붙여서 data를 해싱합니다.
func Encrypt256(data, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(data + salt))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
