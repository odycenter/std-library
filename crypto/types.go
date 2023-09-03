package crypto

import "crypto"

type PaddingType int

// 补位类型
const (
	NONE = PaddingType(iota)
	PaddingZero
	PaddingPKCS5
	PaddingPKCS7
	NoPadding
)

type EncryptType int

// 加密模式
const (
	ECB = EncryptType(iota)
	CBC
	CTR
	OFB
	CFB
)

// 公司密钥对长度
const (
	Bit64   = 64
	Bit128  = 128
	Bit256  = 256
	Bit512  = 512
	Bit1024 = 1024
	Bit2048 = 2048
	Bit4096 = 4096
)

// 私钥类型
const (
	PirTypEC = iota
	PirTypPKCS1
	PirTypPKCS8
)

// 签名方式
const (
	SignTypSHA1   = crypto.SHA1
	SignTypSHA256 = crypto.SHA256
	SignTypSHA512 = crypto.SHA512
	SignTypMD5    = crypto.MD5
)

// 输入参数处理方式
const (
	InputOrigin = iota //默认,输入为加密之后的原型bytes
	InputBase64        //输入数据为base64 bytes
	InputHex           //输入数据为Hex bytes
)
