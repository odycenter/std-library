package sha1_test

import (
	"fmt"
	"github.com/odycenter/std-library/crypto/sha/sha1"
	"testing"
)

func TestSha1(t *testing.T) {
	fmt.Println(sha1.Sum([]byte("a")).Hex())
}
