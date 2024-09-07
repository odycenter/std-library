package aes_test

import (
	"fmt"
	"github.com/odycenter/std-library/crash"
	"github.com/odycenter/std-library/crypto"
	"github.com/odycenter/std-library/crypto/aes"
	"testing"
)

func TestAes_Encrypt(t *testing.T) {
	o := aes.New("aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb", crypto.PaddingPKCS7, crypto.CFB)
	fmt.Println(o.Encrypt("欢迎使用library ").Hex())
}

func TestAes_Decrypt(t *testing.T) {
	o := aes.New("aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb", crypto.PaddingPKCS7, crypto.CFB).WithHex()
	crash.Try(func() {
		//fmt.Println(o.Decrypt([]byte("8d42c12b81c117655380908e36a456550a6d70742ab1de0bc1f9e86f9615447c")).String()) //CBC PKCS7
		//fmt.Println(o.Decrypt([]byte("656565e060f3db4b7b1d23dcb3e4790d2bbaefb2f59cc78c0f21038979870389")).String()) //ECB PKCS7
		//fmt.Println(o.Decrypt([]byte("ffe5174ccfd9a68a6e60109cf9ca8465a105354bf4be41c875a5b924b8e192b1")).String()) //CTR PKCS7
		//fmt.Println(o.Decrypt([]byte("ffe5174ccfd9a68a6e60109cf9ca8465c0b085268756cdb665284664ab344925")).String()) //OFB NONE
		fmt.Println(o.Decrypt([]byte("ffe5174ccfd9a68a6e60109cf9ca84657467d9c19c15f388851230039fad3d27")).String()) //CFB NONE
	}).Catch(nil, func(err error) {
		fmt.Println(err)
	})
}
