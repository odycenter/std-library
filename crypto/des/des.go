// Package des DES对称加密算法实现
package des

import (
	"crypto/cipher"
	cptDes "crypto/des"
	"github.com/odycenter/std-library/crypto"
)

type Des struct {
	key     []byte
	iv      []byte
	padding crypto.PaddingType
	mode    crypto.EncryptType
	decrypt int
}

// New 创建DES加解密对象
func New(key, iv string, padding crypto.PaddingType, mode crypto.EncryptType) *Des {
	d := &Des{
		key:     []byte(key),
		iv:      []byte(iv),
		padding: padding,
		mode:    mode,
	}
	kl := len(d.key)
	ivl := len(d.iv)
	if kl != 8 && kl != 16 && kl != 24 && kl != 32 {
		panic("Des key length must 8/16/24/32.")
	}
	if ivl != 8 && ivl != 16 && ivl != 24 && ivl != 32 && mode != crypto.ECB {
		panic("Des iv length must 8/16/24/32.")
	}
	return d
}

// WithBase64 使用Base64处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (d *Des) WithBase64() *Des {
	d.decrypt = crypto.InputBase64
	return d
}

// WithHex 使用Hex处理输入数据
// Attention 需在 Decrypt 之前调用否则无效
func (d *Des) WithHex() *Des {
	d.decrypt = crypto.InputHex
	return d
}

func (d *Des) cipher(block cipher.Block, de bool) (cipher.BlockMode, cipher.Stream) {
	if de {
		switch d.mode {
		case crypto.ECB:
			return crypto.NewECBDecrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCDecrypter(block, d.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, d.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, d.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBDecrypter(block, d.iv)
		default:
			return nil, nil
		}
	} else {
		switch d.mode {
		case crypto.ECB:
			return crypto.NewECBEncrypter(block), nil
		case crypto.CBC:
			return cipher.NewCBCEncrypter(block, d.iv), nil
		case crypto.CTR:
			return nil, cipher.NewCTR(block, d.iv)
		case crypto.OFB:
			return nil, cipher.NewOFB(block, d.iv)
		case crypto.CFB:
			return nil, cipher.NewCFBEncrypter(block, d.iv)
		default:
			return nil, nil
		}
	}
}

// Encrypt 加密
func (d *Des) Encrypt(origin string) *Ret {
	b := []byte(origin)
	block, err := cptDes.NewCipher(d.key)
	if err != nil {
		return &Ret{err: err}
	}
	padded := crypto.Padding(b, d.padding, block.BlockSize())
	en := make([]byte, len(padded))
	bm, s := d.cipher(block, false)
	if s != nil {
		s.XORKeyStream(en, padded)
		return &Ret{v: en}
	}
	bm.CryptBlocks(en, padded)
	return &Ret{v: en}
}

// Decrypt 解密
func (d *Des) Decrypt(origin []byte) *Ret {
	origin, err := crypto.From(d.decrypt, origin)
	if err != nil {
		return &Ret{err: err}
	}
	block, err := cptDes.NewCipher(d.key)
	if err != nil {
		return &Ret{err: err}
	}
	de := make([]byte, len(origin))
	bm, s := d.cipher(block, true)
	if s != nil {
		s.XORKeyStream(de, origin)
	} else {
		bm.CryptBlocks(de, origin)
	}
	unPadded := crypto.UnPadding(de, d.padding)
	return &Ret{v: unPadded}
}
