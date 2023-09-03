package crypto

import (
	"bytes"
)

// ZeroPadding 补0方式补位
func ZeroPadding(r []byte, blockSize int) []byte {
	padding := blockSize - len(r)%blockSize
	if padding == blockSize {
		return r
	}
	return append(r, bytes.Repeat([]byte{0}, padding)...)
}

// ZeroUnPadding 补0方式去除补位
func ZeroUnPadding(r []byte) []byte {
	return bytes.TrimFunc(r, func(ru rune) bool {
		return ru == rune(0)
	})
}

// PKCS5Padding PKCS5方式补位
func PKCS5Padding(r []byte, blockSize int) []byte {
	padding := blockSize - len(r)%blockSize
	p := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(r, p...)
}

// PKCS5UnPadding PKCS5方式去除补位
func PKCS5UnPadding(r []byte) []byte {
	length := len(r)
	// 去掉最后一个字节 p 次
	p := int(r[length-1])
	return r[:(length - p)]
}

// PKCS7Padding PKCS7方式补位
func PKCS7Padding(r []byte, blockSize int) []byte {
	n := blockSize - len(r)%blockSize
	p := bytes.Repeat([]byte{byte(n)}, n)
	return append(r, p...)
}

// PKCS7UnPadding PKCS7方式去除补位
func PKCS7UnPadding(r []byte) []byte {
	l := len(r)
	p := int(r[l-1])
	return r[:(l - p)]
}

// Padding 补位
// r 原字符串[]byte
// typ 使用补位方式
// bz 位数
func Padding(r []byte, typ PaddingType, bz int) []byte {
	switch typ {
	case PaddingZero:
		return ZeroPadding(r, bz)
	case PaddingPKCS5:
		return PKCS5Padding(r, bz)
	case PaddingPKCS7:
		return PKCS7Padding(r, bz)
	default:
		return r
	}
}

// UnPadding 去除补位
// r 原字符串[]byte
// typ 使用补位方式
func UnPadding(r []byte, typ PaddingType) []byte {
	switch typ {
	case PaddingZero:
		return ZeroUnPadding(r)
	case PaddingPKCS5:
		return PKCS5UnPadding(r)
	case PaddingPKCS7:
		return PKCS7UnPadding(r)
	default:
		return r
	}
}
