package upload

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
)

func (s *service) newHash() hash.Hash {
	return md5.New()
}

func (s *service) writeHash(hash hash.Hash, bytes []byte) (int, error) {
	return hash.Write(bytes)
}

func (s *service) flush(hash hash.Hash) string {
	toString := hex.EncodeToString(hash.Sum(nil))
	hash.Reset()
	return toString
}
