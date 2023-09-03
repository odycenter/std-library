package des_test

import (
	"fmt"
	"std-library/crypto"
	"std-library/crypto/des"
	"testing"
)

func TestDes_Encrypt(t *testing.T) {
	o := des.New("aaaaaaaa", "bbbbbbbb", crypto.PaddingPKCS7, crypto.CFB)
	fmt.Println(o.Encrypt("欢迎使用library ").Hex())
}

func TestDes_Decrypt(t *testing.T) {
	o := des.New("aaaaaaaa", "bbbbbbbb", crypto.PaddingPKCS7, crypto.CFB).WithHex()
	//fmt.Println(o.Decrypt([]byte("9c21ec39e3a5a3fec3763c394817b599fc1caf0d699ad364")).String()) //ECB PKCS7
	//fmt.Println(o.Decrypt([]byte("a5460f99c3863798f8a39dd1faf5bf0bdcc5d1a12bc8bda3")).String()) //CBC PKCS7
	//fmt.Println(o.Decrypt([]byte("048924980e310d6a1bb0d523baf8ba5063db3f245105cf23")).String()) //CTR PKCS7
	//fmt.Println(o.Decrypt([]byte("048924980e310d6a2a51f4a4246e14f19236e9d52adef265")).String()) //OFB PKCS7
	fmt.Println(o.Decrypt([]byte("048924980e310d6a017c9396ccff25f4e3afa30439207636")).String()) //CFB PKCS7
}
