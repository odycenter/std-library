package util_test

import (
	"std-library/app/log/util"
	"testing"
	"time"
)

func BenchmarkIDGenerator(b *testing.B) {
	idGenerator := util.GetIDGenerator()
	for i := 0; i < b.N; i++ {
		idGenerator.Next(time.Now())
	}
}