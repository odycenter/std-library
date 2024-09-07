package array_test

import (
	"fmt"
	"github.com/odycenter/std-library/array"
	"testing"
)

func TestIn(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 1, 1, 2, 3}
	fmt.Println(array.In(s1, 0))
	fmt.Println(array.All(s1, 12))
	fmt.Println(array.Index(s1, 0))
	fmt.Println(array.Last(s1, 0))
	s2 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "0", "0", "1", "1", "2", "3"}
	fmt.Println(array.In(s2, "0"))
	fmt.Println(array.All(s2, "12"))
	fmt.Println(array.Index(s2, "0"))
	fmt.Println(array.Last(s2, "0"))
}
