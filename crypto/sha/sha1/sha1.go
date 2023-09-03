// Package sha1 SHA1摘要生成算法
package sha1

import (
	"crypto/sha1"
)

// Sum 摘要方式Sha1，结果长度40位
func Sum(r []byte) *Ret {
	h := sha1.New()
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}
