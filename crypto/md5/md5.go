// Package md5 MD5摘要生成算法
package md5

import "crypto/md5"

// Sum MD5摘要方式，结果长度32位
func Sum(r []byte) *Ret {
	h := md5.New()
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}
