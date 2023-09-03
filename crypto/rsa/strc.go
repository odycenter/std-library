package rsa

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type Ret struct {
	v   []byte
	err error
}

func (r *Ret) Byte() []byte {
	return r.v
}

// String Attention 不是转换到 Bse64 或 Hex,加密后会输出乱码,用于解密输出字符串
func (r *Ret) String() string {
	return string(r.v)
}

// Hex 转换到Hex十六进制
func (r *Ret) Hex() string {
	return hex.EncodeToString(r.v)
}

// Base64 转换到Base64
func (r *Ret) Base64() string {
	return base64.StdEncoding.EncodeToString(r.v)
}

func (r *Ret) Error() error {
	return r.err
}

func (r *Ret) Result() ([]byte, error) {
	return r.v, r.err
}

func (r *Ret) Bool() bool {
	return r.err == nil
}

type KF struct {
	formatted []byte
}

func (f *KF) Private() string {
	return fmt.Sprintf(`-----BEGIN PUBLIC KEY-----%s-----END PUBLIC KEY-----`, f.formatted)
}

func (f *KF) Public() string {
	return fmt.Sprintf(`-----BEGIN RSA PRIVATE KEY-----%s-----END RSA PRIVATE KEY-----`, f.formatted)
}
