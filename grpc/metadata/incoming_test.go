package metadata_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/grpc/metadata"
)

func TestGetValues(t *testing.T) {
	md := metadata.Pairs("a", "1", "b", "2", "c", "abcdefg", "d", "中")
	ctx, cancel := metadata.NewIncoming().WithCancel().Ctx(md)
	defer cancel()
	fmt.Println(metadata.Get(ctx, "a"))
	fmt.Println(metadata.Get(ctx, "b"))
	fmt.Println(metadata.Get(ctx, "c"))
	fmt.Println(metadata.Get(ctx, "d"))
}
