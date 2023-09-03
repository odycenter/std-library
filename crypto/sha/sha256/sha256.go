// Package sha256 SHA256摘要生成算法
package sha256

import (
	"crypto/sha256"
)

// Sum 摘要方式SHA256，结果长度64
func Sum(r []byte) *Ret {
	h := sha256.New()
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}
