package md5_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/crypto/md5"
)

func TestMd5(t *testing.T) {
	fmt.Println(md5.Sum([]byte("a")).Hex())
	fmt.Println(md5.Sum([]byte("a")).Hex16())
}
