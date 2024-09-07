package bloom_test

import (
	"fmt"
	"github.com/odycenter/std-library/redis/bloom"
	"testing"
)

func TestBloom(t *testing.T) {
	bloom.Init(100000, 0.01)
	bloom.Set([]byte("A"))
	bloom.Set([]byte("B"))
	bloom.Set([]byte("C"))
	bloom.Set([]byte("D"))
	bloom.Set([]byte("E"))
	bloom.Set([]byte("A"))
	bloom.Set([]byte("A"))
	bloom.Set([]byte("A"))
	bloom.Set(bloom.Uint32(1))
	bloom.Set(bloom.Uint32(2))
	bloom.Set([]byte("A"))
	fmt.Println(bloom.Have([]byte("A")))
	fmt.Println(bloom.Have([]byte("B")))
	fmt.Println(bloom.Have([]byte("G")))
	fmt.Println(bloom.Have([]byte("Z")))
	fmt.Println(bloom.Have([]byte("sdwefgw")))
	fmt.Println(bloom.Have(bloom.Uint32(1)))
	fmt.Println(bloom.Have(bloom.Uint32(3)))
}
