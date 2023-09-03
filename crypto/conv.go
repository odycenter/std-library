package crypto

import (
	"encoding/base64"
	"encoding/hex"
)

// Base64 原字符串[]byte转Base64
func Base64(raw []byte) string {
	return base64.StdEncoding.EncodeToString(raw)
}

// FromBase64 Base64转原字符串[]bytes
func FromBase64(raw string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(raw)
}

// Hex 原字符串[]byte转HEX
func Hex(raw []byte) string {
	return hex.EncodeToString(raw)
}

// FromHex HEX转原字符串[]byte
func FromHex(raw string) ([]byte, error) {
	return hex.DecodeString(raw)
}

// From Decrypt数据预处理
func From(t int, raw []byte) ([]byte, error) {
	switch t {
	case InputOrigin:
		return raw, nil
	case InputBase64:
		return FromBase64(string(raw))
	case InputHex:
		return FromHex(string(raw))
	}
	return nil, nil
}
