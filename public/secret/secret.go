package secret

import (
	"crypto/sha256"
	"encoding/hex"
)

func PasswordSum256(pwd string) string {
	sum256 := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(sum256[:])
}
