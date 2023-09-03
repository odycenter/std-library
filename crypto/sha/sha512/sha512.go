// Package sha512 SHA512摘要生成算法
package sha512

import (
	"crypto/sha256"
)

// Sum 摘要方式SHA512，结果长度128
func Sum(r []byte) *Ret {
	h := sha256.New()
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}
