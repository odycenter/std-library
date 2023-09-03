// Package hmac HMAC 摘要生成算法
package hmac

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// Sum1 摘要方式HMACWithSHA1，结果长度40
func Sum1(k, r []byte) *Ret {
	h := hmac.New(sha1.New, k)
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}

// Sum5 摘要方式HMACWithMD5，结果长度32
// 依照标准该方法不推荐使用
func Sum5(k, r []byte) *Ret {
	h := hmac.New(md5.New, k)
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}

// Sum256 摘要方式HMACWithSHA256，结果长度64
func Sum256(k, r []byte) *Ret {
	h := hmac.New(sha256.New, k)
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}

// Sum512 摘要方式HMACWithSHA512，结果长度128
func Sum512(k, r []byte) *Ret {
	h := hmac.New(sha512.New, k)
	h.Write(r)
	return &Ret{v: h.Sum(nil)}
}
