package slice

type T interface {
	any
}

type Inf[V T] interface {
	new() Inf[V]
	Push(v V)
	Index(i int) (V, bool)
	Delete(i int)
	Len() int
	Range(fn func(v V) error) (err error)
}

// New 创建指定的map
func New[V T](i Inf[V]) Inf[V] {
	return i.new()
}

// ArrayPartition
/*
範例1：
陣列：[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]，正整數：2
期望結果: [[1, 2], [3, 4], [5, 6], [7, 8], [9, 10]]
呼叫: res:= ArrayPartition(arr,2)
*/
func ArrayPartition[T any](arr []T, eachPartitionSize int64) [][]T {
	length := int64(len(arr))
	if eachPartitionSize < 1 || eachPartitionSize >= length {
		return [][]T{arr}
	}

	//取得應該數組分割為多少份
	var quantity int64
	if length%eachPartitionSize == 0 {
		quantity = length / eachPartitionSize
	} else {
		quantity = (length / eachPartitionSize) + 1
	}

	partitions := make([][]T, quantity)
	start, end := int64(0), eachPartitionSize
	for i := int64(0); i < quantity; i++ {
		if end > length {
			end = length
		}
		partitions[i] = arr[start:end]
		start, end = end, end+eachPartitionSize
	}

	return partitions
}

// ToMap
// 將 slice 轉成 map
func ToMap[K comparable, V any](arr []V, keyFunc func(vl V) K) map[K]V {
	m := make(map[K]V, len(arr))
	for _, v := range arr {
		m[keyFunc(v)] = v
	}
	return m
}
