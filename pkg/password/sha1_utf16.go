package password

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

type Hasher interface {
	Hash(plain string) (string, error)
}

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Hash(plain string) (string, error) {
	encoder := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewEncoder()
	b, err := encoder.Bytes([]byte(plain))
	if err != nil {
		return "", err
	}

	digest := sha1.New()
	digest.Write(b)
	sum := digest.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(sum)), nil
}
