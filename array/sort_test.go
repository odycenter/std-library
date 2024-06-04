package array_test

import (
	"fmt"
	"std-library/array"
	"testing"
)

func TestSort(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5, 7, 8, 9, 5, 3, 2}
	array.Sort(a)
	fmt.Println(a)
	fmt.Println(array.IsAsc(a))
	var b = []string{"cat", "Door", "Apple", "banana"}
	array.Sort(b)
	fmt.Println(b)
	fmt.Println(array.IsAsc(b))
}

func TestReverse(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5, 7, 8, 9, 5, 3, 2}
	array.SortReverse(a)
	fmt.Println(a)
	fmt.Println(array.IsDesc(a))
	var b = []string{"cat", "Door", "Apple", "banana"}
	array.SortReverse(b)
	fmt.Println(b)
	fmt.Println(array.IsDesc(b))
}

func TestSortFloat(t *testing.T) {
	var a = []float64{1.1, 2.3, 3.141, 4.564, 5.65, 7.99, 8.34, 9.23, 5.44, 3.22, 2.67}
	array.SortFloat(a)
	fmt.Println(a)
	fmt.Println(array.IsDesc(a))
}
