package array

import (
	"reflect"
	"sort"
	"unsafe"
)

type Asc[T Numeric | string] []T
type Desc[T Numeric | string] []T

func (inf Asc[T]) Len() int           { return len(inf) }
func (inf Asc[T]) Less(i, j int) bool { return inf[i] < inf[j] }
func (inf Asc[T]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

func (inf Desc[T]) Len() int           { return len(inf) }
func (inf Desc[T]) Less(i, j int) bool { return inf[i] > inf[j] }
func (inf Desc[T]) Swap(i, j int)      { inf[i], inf[j] = inf[j], inf[i] }

// Sort 正序排序数组-由小到大
func Sort[T Numeric | string](l []T) {
	sort.Sort(Asc[T](l))
}

// SortReverse 倒序序排序数组-由大到小
func SortReverse[T Numeric | string](l []T) {
	sort.Sort(Desc[T](l))
}

// IsAsc 是否已经正序
func IsAsc[T Numeric | string](l []T) bool {
	return sort.IsSorted(Asc[T](l))
}

// IsDesc 是否已经倒序
func IsDesc[T Numeric | string](l []T) bool {
	return sort.IsSorted(Desc[T](l))
}

// SortFloat float专用快速排序-正序，如果是超大型的float数组，效率优先情况下应使用该方法
func SortFloat[T float64 | float32](l []T) {
	var c []int
	aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&l))
	cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*cHdr = *aHdr
	sort.Ints(c)
}
