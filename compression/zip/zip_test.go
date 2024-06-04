package zip

import (
	"testing"
)

func TestZip(t *testing.T) {
	t.Log(Compress("test", "test.zip", "12345678"))
	t.Log(Compress("test.zip", "test1.zip", "87654321"))
}
