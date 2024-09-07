// Package ecdsa ECDSA加密算法
package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/odycenter/std-library/crypto"
	"github.com/odycenter/std-library/crypto/sha/sha1"
	"github.com/odycenter/std-library/crypto/sha/sha256"
	"github.com/odycenter/std-library/crypto/sha/sha512"
	"math/big"
)

type ECDSA struct {
	curve   curve
	hash    hash
	decrypt int
}

// New 创建 ecdsa 签名对象
// curve 椭圆曲线生成方式
// hash 原文摘要方式
func New(curve curve, hash hash) *ECDSA {
	return &ECDSA{
		curve: curve,
		hash:  hash,
	}
}

// WithBase64 使用Base64处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (c *ECDSA) WithBase64() *ECDSA {
	c.decrypt = crypto.InputBase64
	return c
}

// WithHex 使用Hex处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (c *ECDSA) WithHex() *ECDSA {
	c.decrypt = crypto.InputHex
	return c
}

type curve int

// 椭圆曲线生成方式
const (
	P224 = curve(iota)
	P256
	P384
	P521
)

func ellipticCurve(t curve) elliptic.Curve {
	switch t {
	case P224:
		return elliptic.P224()
	case P256:
		return elliptic.P256()
	case P384:
		return elliptic.P384()
	case P521:
		return elliptic.P521()
	default:
		return nil
	}
}

type hash int

// 原文摘要方式
const (
	SHA1 = hash(iota)
	SHA256
	SHA512
)

func genHash(d []byte, t hash) []byte {
	switch t {
	case SHA1:
		return sha1.Sum(d).Byte()
	case SHA256:
		return sha256.Sum(d).Byte()
	case SHA512:
		return sha512.Sum(d).Byte()
	default:
		return nil
	}
}

// GenKeys 生成 ecdsa 公私密钥对
func (c *ECDSA) GenKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(ellipticCurve(c.curve), rand.Reader)
	if err != nil {
		return nil, nil
	}
	return privateKey, &privateKey.PublicKey
}

var ErrArgs = errors.New("invalid args")

// PrivateEncode 将private key转为Pem输出的目标格式
func PrivateEncode(priKey *ecdsa.PrivateKey) *Ret {
	if priKey == nil {
		return &Ret{err: ErrArgs}
	}
	ecPri, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return &Ret{err: ErrArgs}
	}
	block := pem.Block{
		Type:    "ECDSA PRIVATE KEY",
		Headers: nil,
		Bytes:   ecPri,
	}
	return &Ret{
		v:   pem.EncodeToMemory(&block),
		err: nil,
	}
}

// PublicEncode 将public key转为Pem输出的目标格式
func PublicEncode(pubKey *ecdsa.PublicKey) *Ret {
	if pubKey == nil {
		return &Ret{err: ErrArgs}
	}
	ecPub, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return &Ret{err: ErrArgs}
	}
	block := pem.Block{
		Type:    "ECDSA PUBLIC KEY",
		Headers: nil,
		Bytes:   ecPub,
	}
	return &Ret{
		v:   pem.EncodeToMemory(&block),
		err: nil,
	}
}

// PrivateDecode 将Pem []byte的private key转为ecdsa.PrivateKey
func PrivateDecode(priKey []byte) (pri *ecdsa.PrivateKey, err error) {
	block, _ := pem.Decode(priKey)
	return x509.ParseECPrivateKey(block.Bytes)
}

// PublicDecode 将Pem []byte的public key转为ecdsa.PublicKey
func PublicDecode(pubKey []byte) (pub *ecdsa.PublicKey, err error) {
	block, _ := pem.Decode(pubKey)
	v, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return v.(*ecdsa.PublicKey), err
}

// Sign 签名
func (c *ECDSA) Sign(origin []byte, key *ecdsa.PrivateKey) (retR *Ret, retS *Ret) {
	h := genHash(origin, c.hash)
	r, s, err := ecdsa.Sign(rand.Reader, key, h[:])
	if err != nil {
		retR = &Ret{nil, err}
		retR = &Ret{nil, err}
		return
	}
	rs, err := r.MarshalText()
	retR = &Ret{rs, err}
	ss, err := s.MarshalText()
	retS = &Ret{ss, err}
	return
}

// BigInt 将[]byte转为big.Int
func BigInt(b []byte) *big.Int {
	i := big.Int{}
	err := i.UnmarshalText(b)
	if err != nil {
		return nil
	}
	return &i
}

// Verify 验证签名
func (c *ECDSA) Verify(origin []byte, signR, signS []byte, key *ecdsa.PublicKey) bool {
	signR, err := crypto.From(c.decrypt, signR)
	if err != nil {
		return false
	}
	signS, err = crypto.From(c.decrypt, signS)
	if err != nil {
		return false
	}
	r, s := BigInt(signR), BigInt(signS)
	h := genHash(origin, c.hash)
	return ecdsa.Verify(key, h, r, s)
}

// SignASN1 以ASN1方式进行签名(推荐)
func (c *ECDSA) SignASN1(origin []byte, key *ecdsa.PrivateKey) *Ret {
	h := genHash(origin, c.hash)
	sig, err := ecdsa.SignASN1(rand.Reader, key, h[:])
	return &Ret{sig, err}
}

// VerifyASN1 以ASN1方式进行验证签名(推荐)
func (c *ECDSA) VerifyASN1(origin, sig []byte, key *ecdsa.PublicKey) bool {
	sig, err := crypto.From(c.decrypt, sig)
	if err != nil {
		return false
	}
	h := genHash(origin, c.hash)
	return ecdsa.VerifyASN1(key, h, sig)
}
