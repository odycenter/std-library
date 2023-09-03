// Package rsa RSA非对称加密算法及密钥生成
package rsa

import (
	"bytes"
	cpt "crypto"
	"crypto/md5"
	"crypto/rand"
	cptRsa "crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
	"std-library/crypto"
)

type Rsa struct {
	pubKey  []byte
	priKey  []byte
	priTyp  int
	signTyp cpt.Hash
	decrypt int
}

// New 创建RSA加解密对象
func New(pubKey, priKey string, signTyp ...cpt.Hash) *Rsa {
	r := &Rsa{
		pubKey: []byte(pubKey),
		priKey: []byte(priKey),
	}
	if len(signTyp) > 0 {
		r.signTyp = signTyp[0]
	}
	if priKey == "" {
		return r
	}
	err := r.checkPri()
	if err != nil {
		panic(err.Error())
	}
	return r
}

// WithBase64 使用Base64处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (r *Rsa) WithBase64() *Rsa {
	r.decrypt = crypto.InputBase64
	return r
}

// WithHex 使用Hex处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (r *Rsa) WithHex() *Rsa {
	r.decrypt = crypto.InputHex
	return r
}

// WithPrivateParse 使用指定的私钥处理方法
// crypto.PirTypEC 暂不提供
// crypto.PirTypPKCS1 PKCS1方式的私钥
// crypto.PirTypPKCS8 PKCS8方式的私钥
func (r *Rsa) WithPrivateParse(priParseTyp int) *Rsa {
	r.priTyp = priParseTyp
	return r
}

func (r *Rsa) checkPri() error {
	block, _ := pem.Decode(r.priKey)
	if _, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		r.priTyp = crypto.PirTypEC
		return nil
	}
	if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		r.priTyp = crypto.PirTypPKCS1
		return nil
	}
	if _, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		r.priTyp = crypto.PirTypPKCS8
		return nil
	}
	return errors.New("unknown private key Type")
}

func (r *Rsa) parsePriKey(block *pem.Block) (pri *cptRsa.PrivateKey, err error) {
	switch r.priTyp {
	case crypto.PirTypEC:
		if pri, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return pri, nil
		}
	case crypto.PirTypPKCS1:
		if pri, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return pri, nil
		}
	case crypto.PirTypPKCS8:
		priIface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			return priIface.(*cptRsa.PrivateKey), nil
		}
	}
	return nil, err
}

func (r *Rsa) hash() hash.Hash {
	switch r.signTyp {
	case crypto.SignTypSHA1:
		return sha1.New()
	case crypto.SignTypSHA256:
		return sha256.New()
	case crypto.SignTypSHA512:
		return sha512.New()
	case crypto.SignTypMD5:
		return md5.New()
	}
	return nil
}

// Encrypt 加密
func (r *Rsa) Encrypt(origin []byte) *Ret {
	block, _ := pem.Decode(r.pubKey)
	if block == nil {
		return &Ret{err: errors.New("Rsa:Public Key Error")}
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return &Ret{err: err}
	}
	result, err := cptRsa.EncryptPKCS1v15(rand.Reader, pub.(*cptRsa.PublicKey), origin)
	return &Ret{result, err}
}

// Decrypt 解密
func (r *Rsa) Decrypt(origin []byte) *Ret {
	origin, err := crypto.From(r.decrypt, origin)
	if err != nil {
		return &Ret{err: err}
	}
	block, _ := pem.Decode(r.priKey)
	if block == nil {
		return &Ret{err: errors.New("Rsa:Private Key Error")}
	}
	pri, err := r.parsePriKey(block)
	if err != nil {
		return &Ret{err: err}
	}
	bs := splitWithSize(origin, (pri.N.BitLen()+7)/8)
	var buf bytes.Buffer
	for _, b := range bs {
		de, err := cptRsa.DecryptPKCS1v15(rand.Reader, pri, b)
		if err != nil {
			return &Ret{err: err}
		}
		buf.Write(de)
	}
	return &Ret{v: buf.Bytes(), err: nil}
}

// Sign 签名
func (r *Rsa) Sign(origin []byte) *Ret {
	h := r.hash()
	h.Write(origin)
	sum := h.Sum(nil)
	block, _ := pem.Decode(r.priKey)
	if block == nil {
		return &Ret{err: errors.New("Rsa:Private Key Error")}
	}
	pri, err := r.parsePriKey(block)
	if err != nil {
		return &Ret{err: err}
	}
	sign, err := cptRsa.SignPKCS1v15(rand.Reader, pri, r.signTyp, sum)

	return &Ret{v: sign, err: err}
}

// Verify 验证签名
func (r *Rsa) Verify(origin, sign []byte) *Ret {
	sign, err := crypto.From(r.decrypt, sign)
	if err != nil {
		return &Ret{err: err}
	}
	h := r.hash()
	h.Write(origin)
	sum := h.Sum(nil)
	block, _ := pem.Decode(r.pubKey)
	if block == nil {
		return &Ret{err: errors.New("Rsa:Public Key Error")}
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return &Ret{err: err}
	}
	return &Ret{err: cptRsa.VerifyPKCS1v15(pub.(*cptRsa.PublicKey), r.signTyp, sum, sign)}
}

// 根据长度切割字节
func splitWithSize(plain []byte, size int) [][]byte {
	var result [][]byte
	plainLen := len(plain)
	for i := 0; i < plainLen/size; i++ {
		result = append(result, plain[size*i:size*(i+1)])
	}
	plainMod := plainLen % size
	if plainMod > 0 {
		result = append(result, plain[plainLen-plainMod:])
	}
	return result
}
