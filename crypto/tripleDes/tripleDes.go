// Package tripleDes 三重DES对称加密算法
package tripleDes

import (
	"crypto/cipher"
	"crypto/des"
	"std-library/crypto"
)

type Tripledes struct {
	key     []byte
	iv      []byte
	padding crypto.PaddingType
	mode    crypto.EncryptType
	decrypt int
}

// New 创建3DES加解密对象
func New(key, iv string, padding crypto.PaddingType, mode crypto.EncryptType) *Tripledes {
	d := &Tripledes{
		key:     []byte(key),
		iv:      []byte(iv),
		padding: padding,
		mode:    mode,
	}
	kl := len(d.key)
	ivl := len(d.iv)
	if kl != 24 && kl != 32 {
		panic("Tripledes key length must 24/32.")
	}
	if ivl != 8 && ivl != 16 && ivl != 24 && ivl != 32 && mode != crypto.ECB {
		panic("Tripledes iv length must 8/16/24/32.")
	}
	return d
}

// WithBase64 使用Base64处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (t *Tripledes) WithBase64() *Tripledes {
	t.decrypt = crypto.InputBase64
	return t
}

// WithHex 使用Hex处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (t *Tripledes) WithHex() *Tripledes {
	t.decrypt = crypto.InputHex
	return t
}
func (t *Tripledes) cipher(block cipher.Block, de bool) (cipher.BlockMode, cipher.Stream) {
	if de {
		switch t.mode {
		case crypto.ECB:
			return crypto.NewECBDecrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCDecrypter(block, t.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, t.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, t.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBDecrypter(block, t.iv)
		default:
			return nil, nil
		}
	} else {
		switch t.mode {
		case crypto.ECB:
			return crypto.NewECBEncrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCEncrypter(block, t.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, t.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, t.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBEncrypter(block, t.iv)
		default:
			return nil, nil
		}
	}
}

// Encrypt 加密
func (t *Tripledes) Encrypt(origin string) *Ret {
	b := []byte(origin)
	block, err := des.NewTripleDESCipher(t.key)
	if err != nil {
		return &Ret{err: err}
	}
	padded := crypto.Padding(b, t.padding, block.BlockSize())
	en := make([]byte, len(padded))
	bm, s := t.cipher(block, false)
	if s != nil {
		s.XORKeyStream(en, padded)
		return &Ret{v: en}
	}
	bm.CryptBlocks(en, padded)
	return &Ret{v: en}
}

// Decrypt 解密
func (t *Tripledes) Decrypt(origin []byte) *Ret {
	origin, err := crypto.From(t.decrypt, origin)
	if err != nil {
		return &Ret{err: err}
	}
	block, err := des.NewTripleDESCipher(t.key)
	if err != nil {
		return &Ret{err: err}
	}
	de := make([]byte, len(origin))
	bm, s := t.cipher(block, true)
	if s != nil {
		s.XORKeyStream(de, origin)
	} else {
		bm.CryptBlocks(de, origin)
	}
	unPadded := crypto.UnPadding(de, t.padding)
	return &Ret{v: unPadded}
}
