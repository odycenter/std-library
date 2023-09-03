package operator_test

import (
	"fmt"
	"std-library/operator"
	"testing"
)

func TestIF(t *testing.T) {
	fmt.Println(operator.IF(1 == 1, "1", "2"))
	fmt.Println(operator.IF(1 == 2, "1", "2"))
	fmt.Println(operator.IF(1 == 2, true, false))
	fmt.Println(operator.IF(1 == 2, 123, 456))
	fmt.Println(operator.IF(1 == 2, 1.23, 4.56))
	type s struct {
		A string
	}
	var s1 = s{A: "a"}
	var s2 = s{A: "b"}
	fmt.Printf("%#v", operator.IF(1 == 2, s1, s2))
}
