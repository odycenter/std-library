package rsa

import (
	"bytes"
	"crypto/rand"
	cptRsa "crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
)

// KeyPair 获取密钥对
func KeyPair(bits int) (pubPem, priPem *Ret, err error) {
	priKey, err := cptRsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	priKeyDer := x509.MarshalPKCS1PrivateKey(priKey)
	priKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priKeyDer,
	}
	priKeyPem := pem.EncodeToMemory(&priKeyBlock)
	pubKey := priKey.PublicKey
	pubKeyDer, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return nil, nil, err
	}
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKeyDer,
	}
	pubKeyPem := pem.EncodeToMemory(&pubKeyBlock)
	return &Ret{v: pubKeyPem}, &Ret{v: priKeyPem}, nil
}

func KeyFormat(key string) *KF {
	slice := strings.Split(key, "-----")
	// 找出公钥串或者密钥串
	var k string

	for _, v := range slice {
		if len(v) > 21 { // 因为这个 `BEGIN RSA PRIVATE KEY` 长度为21
			k = v
			break
		}
	}

	// 找不到公钥/密钥
	if len(k) == 0 {
		return nil
	}

	k = strings.ReplaceAll(strings.ReplaceAll(k, "\n", ""), "\r", "")
	keyLen := len(k)
	buf := bytes.Buffer{}
	buf.WriteByte('\n')
	offset := 0
	for {

		if end := offset + 64; end > keyLen {
			buf.WriteString(k[offset:])
			buf.WriteByte('\n')
			break
		} else {
			buf.WriteString(k[offset:end])
			buf.WriteByte('\n')
			offset = end
		}
	}
	return &KF{buf.Bytes()}
}
