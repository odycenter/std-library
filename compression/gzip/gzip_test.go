package gzip_test

import (
	"encoding/base64"
	"fmt"
	"std-library/compression/gzip"
	"testing"
)

func TestGzip(t *testing.T) {
	o := []byte("压缩压缩文件Gzip，测试测试")
	b, err := gzip.Compress(o)
	if err != nil {
		panic(err)
	}
	fmt.Println("len o:", len(o), "len b:", len(b))
	fmt.Printf("%s\n", b)
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(b))
	b, err = gzip.Decompress(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
