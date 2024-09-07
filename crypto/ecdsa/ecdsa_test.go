package ecdsa_test

import (
	"encoding/base64"
	"fmt"
	"github.com/odycenter/std-library/crypto/ecdsa"
	"testing"
)

func TestNew(t *testing.T) {
	b := []byte("欢迎使用library ")
	o := ecdsa.New(ecdsa.P224, ecdsa.SHA1).WithBase64()
	pri, pub := o.GenKeys()
	pr := ecdsa.PrivateEncode(pri)
	pu := ecdsa.PublicEncode(pub)
	fmt.Println("keys:\n", pr.Base64(), "\n", pu.Base64())
	ret := o.SignASN1(b, pri).Base64()
	fmt.Println("sign:", ret)
	pb, _ := base64.StdEncoding.DecodeString(pu.Base64())
	pub, _ = ecdsa.PublicDecode(pb)
	fmt.Println(o.VerifyASN1(b, []byte(ret), pub))
	pri, pub = o.GenKeys()
	r, s := o.Sign(b, pri)
	rb, sb := r.Base64(), s.Base64()
	fmt.Println("sign:", rb, sb)
	fmt.Println(o.Verify(b, []byte(rb), []byte(sb), pub))
}
