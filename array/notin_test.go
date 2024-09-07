package array_test

import (
	"fmt"
	"github.com/odycenter/std-library/array"
	"testing"
)

func TestNot(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 1, 1, 2, 3}
	s2 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "0", "0", "1", "1", "2", "3"}
	fmt.Println(array.Not(s1, 1))
	fmt.Println(array.Not(s1, 22))
	fmt.Println(array.NoOne(s2, []string{"1", "2"}))
	fmt.Println(array.NoOne(s2, []string{"22", "33"}))
}
