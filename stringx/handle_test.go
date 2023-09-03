package stringx_test

import (
	"fmt"
	"std-library/array"
	"std-library/stringx"
	"testing"
)

func TestTrimHtml(t *testing.T) {
	fmt.Println(stringx.TrimHtml(`剩余部分1<div class="Popover js-hovercard-content position-absolute" style="display: none; outline: none;" tabindex="0">
  <div class="Popover-message Popover-message--bottom-left Popover-message--large Box color-shadow-large" style="width:360px;"></div>
</div>剩余部分2`))
}
func TestSplit(t *testing.T) {
	s1 := "1,2,3,4,5"
	s2 := "a,b,c,d,e"
	fmt.Printf("%#v\n", stringx.Split(s1, ",", false).Ints())
	fmt.Printf("%#v\n", stringx.Split(s2, ",", false).Strings())
}
func TestJoin(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5}
	a2 := []string{"a", "b", "c", "d", "e"}
	fmt.Println(stringx.Join(a1, ","))
	fmt.Println(stringx.Join(a2, ","))
}

func TestHidden(t *testing.T) {
	fmt.Println(stringx.Hidden("Hello World!", 3, 15))
}

func TestSub(t *testing.T) {
	fmt.Println(stringx.Sub("Hello World!", 2, 15))
}

func TestTrim(t *testing.T) {
	fmt.Println(stringx.Trim("$123,456,789.00", func(r rune) bool {
		return r == '$' || r == ',' || r == '.'
	}))
}

func TestCamel2Snack(t *testing.T) {
	fmt.Println(stringx.Camel2Snake("HelloWorld1234@@"))
	fmt.Println(stringx.Camel2Snake("helloWorld1234@@"))
	fmt.Println(stringx.Camel2Snake("message"))
	fmt.Println(stringx.Camel2Snake("中文English1234"))
}

func TestRemoveDuplicates(t *testing.T) {
	s := "a,b,c,d,a,b,c,d"
	fmt.Println(stringx.RemoveDuplicates(s))
	fmt.Println(array.RemoveDuplicates([]int{1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 1, 2}))
	fmt.Println(array.RemoveDuplicates(stringx.Split(s, ",").Strings()))
	fmt.Println(stringx.Join(array.RemoveDuplicates(stringx.Split(s, ",").Strings()), ","))
}
