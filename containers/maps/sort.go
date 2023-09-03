package maps

import "sort"

type Numeric interface {
	uint8 |
		uint16 |
		uint32 |
		uint64 |
		int8 |
		int16 |
		int32 |
		int64 |
		float32 |
		float64 |
		int |
		uint
}

type KV[K, V Numeric | string] struct {
	Key   K
	Value V
}

type KAsc[K, V Numeric | string] []KV[K, V]
type KDesc[K, V Numeric | string] []KV[K, V]
type VAsc[K, V Numeric | string] []KV[K, V]
type VDesc[K, V Numeric | string] []KV[K, V]

func (inf VAsc[K, V]) Len() int           { return len(inf) }
func (inf VAsc[K, V]) Less(i, j int) bool { return inf[i].Value < inf[j].Value }
func (inf VAsc[K, V]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

func (inf VDesc[K, V]) Len() int           { return len(inf) }
func (inf VDesc[K, V]) Less(i, j int) bool { return inf[i].Value > inf[j].Value }
func (inf VDesc[K, V]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

func (inf KAsc[K, V]) Len() int           { return len(inf) }
func (inf KAsc[K, V]) Less(i, j int) bool { return inf[i].Key < inf[j].Key }
func (inf KAsc[K, V]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

func (inf KDesc[K, V]) Len() int           { return len(inf) }
func (inf KDesc[K, V]) Less(i, j int) bool { return inf[i].Key > inf[j].Key }
func (inf KDesc[K, V]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

// SortByVal 根据map的value排序map并返回 KV 结构升序排序数组
func SortByVal[K, V Numeric | string](m map[K]V) VAsc[K, V] {
	kvs := make(VAsc[K, V], len(m))
	var i int
	for k, v := range m {
		kvs[i] = KV[K, V]{k, v}
		i++
	}
	sort.Sort(kvs)
	return kvs
}

// SortByValReverse 根据map的value排序map并返回 KV 结构降序排序数组
func SortByValReverse[K, V Numeric | string](m map[K]V) VDesc[K, V] {
	kvs := make(VDesc[K, V], len(m))
	var i int
	for k, v := range m {
		kvs[i] = KV[K, V]{k, v}
		i++
	}
	sort.Sort(kvs)
	return kvs
}

// SortByKey 根据map的key排序map并返回 KV 结构升序排序数组
func SortByKey[K, V Numeric | string](m map[K]V) KAsc[K, V] {
	kvs := make(KAsc[K, V], len(m))
	var i int
	for k, v := range m {
		kvs[i] = KV[K, V]{k, v}
		i++
	}
	sort.Sort(kvs)
	return kvs
}

// SortByKeyReverse 根据map的key排序map并返回 KV 结构降序排序数组
func SortByKeyReverse[K, V Numeric | string](m map[K]V) KDesc[K, V] {
	kvs := make(KDesc[K, V], len(m))
	var i int
	for k, v := range m {
		kvs[i] = KV[K, V]{k, v}
		i++
	}
	sort.Sort(kvs)
	return kvs
}
