package hash_test

import (
	"github.com/odycenter/std-library/crypto/hash"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSha256Hex(t *testing.T) {
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", hash.Sha256Hex(""))
	assert.Equal(t, "a318c24216defe206feeb73ef5be00033fa9c4a74d0b967f6532a26ca5906d3b", hash.Sha256Hex("+"))
	assert.Equal(t, "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", hash.Sha256Hex("123"))
}

func TestMd5Hex(t *testing.T) {
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", hash.Md5Hex(""))
	assert.Equal(t, "26b17225b626fb9238849fd60eabdf60", hash.Md5Hex("+"))
	assert.Equal(t, "202cb962ac59075b964b07152d234b70", hash.Md5Hex("123"))
}
