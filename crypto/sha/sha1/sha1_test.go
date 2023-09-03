package sha1_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/crypto/sha/sha1"
)

func TestSha1(t *testing.T) {
	fmt.Println(sha1.Sum([]byte("a")).Hex())
}
