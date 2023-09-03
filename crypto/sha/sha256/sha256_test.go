package sha256_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/crypto/sha/sha256"
)

func TestSha256(t *testing.T) {
	fmt.Println(sha256.Sum([]byte("a")).Hex())
}
