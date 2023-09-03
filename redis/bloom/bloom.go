// Package bloom 布隆过滤器
package bloom

import (
	"encoding/binary"
	"github.com/bits-and-blooms/bloom/v3"
)

// Init gen a bloom filter
// n:elements number
// fp:false-positive rate
var filter *bloom.BloomFilter

// Init 初始化布隆过滤器
// 默认创建每1000000条数据保持误差1%的bloom过滤器
func Init(n uint, fp float64) {
	if n == 0 {
		n = 1000000
	}
	if fp == 0 {
		fp = 0.01
	}
	filter = bloom.NewWithEstimates(n, fp)
}

// Set 插入[]byte值
func Set(elem []byte) {
	filter.Add(elem)
}

// Uint32 插入uint32类型数据
func Uint32(elem uint32) []byte {
	n1 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, elem)
	return n1
}

// Have 检查是否存在
// 注意：实际使用中，该值返回true并不能表示查询的数据一定存在于Cache或DB
// 返回false则可以表示数据一定不存在于Cache或DB
func Have(elem []byte) bool {
	return filter.Test(elem)
}

// Clear 清空bloom过滤器
func Clear() {
	filter.ClearAll()
}
