package sha512_test

import (
	"fmt"
	"github.com/odycenter/std-library/crypto/sha/sha256"
	"testing"
)

func TestSha512(t *testing.T) {
	fmt.Println(sha256.Sum([]byte("a")).Hex())
}
