// Package aes AES加密算法
package aes

import (
	cptAes "crypto/aes"
	"crypto/cipher"

	"github.com/odycenter/std-library/crypto"
)

type Aes struct {
	key     []byte
	iv      []byte
	padding crypto.PaddingType
	mode    crypto.EncryptType
	decrypt int
}

// New 创建AES加解密对象，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func New(key, iv string, padding crypto.PaddingType, mode crypto.EncryptType) *Aes {
	a := &Aes{
		key:     []byte(key),
		iv:      []byte(iv),
		padding: padding,
		mode:    mode,
	}
	l := len(a.key)
	if l != 16 && l != 24 && l != 32 {
		panic("AES key length must 16/24/32 Bits")
	}
	return a
}

// WithBase64 使用Base64处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (a *Aes) WithBase64() *Aes {
	a.decrypt = crypto.InputBase64
	return a
}

// WithHex 使用Hex处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (a *Aes) WithHex() *Aes {
	a.decrypt = crypto.InputHex
	return a
}

func (a *Aes) cipher(block cipher.Block, de bool) (cipher.BlockMode, cipher.Stream) {
	if de {
		switch a.mode {
		case crypto.ECB:
			return crypto.NewECBDecrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCDecrypter(block, a.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, a.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, a.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBDecrypter(block, a.iv)
		default:
			return nil, nil
		}
	} else {
		switch a.mode {
		case crypto.ECB:
			return crypto.NewECBEncrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCEncrypter(block, a.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, a.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, a.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBEncrypter(block, a.iv)
		default:
			return nil, nil
		}
	}
}

// Encrypt 加密
func (a *Aes) Encrypt(origin string) *Ret {
	b := []byte(origin)
	block, err := cptAes.NewCipher(a.key)
	if err != nil {
		return &Ret{err: err}
	}
	var padded []byte
	bs := block.BlockSize()
	padded = crypto.Padding(b, a.padding, bs)
	en := make([]byte, len(padded))
	bm, s := a.cipher(block, false)
	if s != nil {
		s.XORKeyStream(en, padded)
		return &Ret{v: en}
	}
	bm.CryptBlocks(en, padded)
	return &Ret{v: en}
}

// Decrypt 解密
func (a *Aes) Decrypt(origin []byte) *Ret {
	origin, err := crypto.From(a.decrypt, origin)
	if err != nil {
		return &Ret{err: err}
	}
	block, err := cptAes.NewCipher(a.key)
	if err != nil {
		return &Ret{err: err}
	}
	de := make([]byte, len(origin))
	bm, s := a.cipher(block, true)
	if s != nil {
		s.XORKeyStream(de, origin)
	} else {
		bm.CryptBlocks(de, origin)
	}
	unPadded := crypto.UnPadding(de, a.padding)
	return &Ret{v: unPadded}
}
