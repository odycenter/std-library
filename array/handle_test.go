package array_test

import (
	"fmt"
	"std-library/array"
	"testing"
)

func TestReplace(t *testing.T) {
	fmt.Printf("%#v\n", array.Replace([]string{"'A'", "'B'", "'大'"}, "'"))
}

func TestReserve(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 5, 4, 3}
	fmt.Println(array.Reverse(a))
}

func TestRemove(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 8, 2, 2, 2, 2, 2, 2}
	fmt.Println(array.Remove(a, 1, 1))
	fmt.Println(array.Remove(a, 2, 2))
	fmt.Println(array.Remove(a, 2, -1))
}

func TestRemoveFn(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Println(array.RemoveFn(a, func(elem int) bool {
		return elem == 5
	}))
}
