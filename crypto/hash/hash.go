package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hex(value string) string {
	return Sha256HexByByte([]byte(value))
}

func Sha256HexByByte(value []byte) string {
	h := sha256.New()
	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}

func Sha1Hex(value string) string {
	return Sha1HexByByte([]byte(value))
}

func Sha1HexByByte(value []byte) string {
	h := sha1.New()
	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}

func Md5Hex(value string) string {
	return Md5HexByByte([]byte(value))
}

func Md5HexByByte(value []byte) string {
	h := md5.New()
	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}
