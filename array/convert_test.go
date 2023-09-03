package array_test

import (
	"fmt"
	"std-library/array"
	"testing"
)

func TestNums2Strs(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(array.Nums2Strings(a))
}

func TestStrs2Nums(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5", "6"}
	fmt.Println(array.Strings2Nums[uint8](a, array.Uint8))
	fmt.Println(array.Strings2Nums[uint16](a, array.Uint16))
	fmt.Println(array.Strings2Nums[uint32](a, array.Uint32))
	fmt.Println(array.Strings2Nums[uint64](a, array.Uint64))
	fmt.Println(array.Strings2Nums[float32](a, array.Float32))
	fmt.Println(array.Strings2Nums[float64](a, array.Float64))
	fmt.Println(array.Strings2Nums[int8](a, array.Int8))
	fmt.Println(array.Strings2Nums[int8](a, array.Int8))
	fmt.Println(array.Strings2Nums[int16](a, array.Int16))
	fmt.Println(array.Strings2Nums[int32](a, array.Int32))
	fmt.Println(array.Strings2Nums[int64](a, array.Int64))
	fmt.Println(array.Strings2Nums[int](a, array.Int))
	fmt.Println(array.Strings2Nums[uint](a, array.Uint))
}

func TestMap(t *testing.T) {
	a1 := []string{"1", "2", "3", "4", "5", "6"}
	fmt.Printf("%#v\n", array.Map(a1))
	a2 := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%#v\n", array.Map(a2))
}

func TestBytes2String(t *testing.T) {
	fmt.Println(array.Bytes2String([]byte{97, 98}))
}
