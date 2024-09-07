package stringx_test

import (
	"fmt"
	"github.com/odycenter/std-library/stringx"
	"testing"
)

func TestStr2Num(t *testing.T) {
	fmt.Println(stringx.Str2Num[uint8]("999999", stringx.Uint8))
	fmt.Println(stringx.Str2Num[uint16]("999999", stringx.Uint16))
	fmt.Println(stringx.Str2Num[uint32]("999999", stringx.Uint32))
	fmt.Println(stringx.Str2Num[uint64]("999999", stringx.Uint64))
	fmt.Println(stringx.Str2Num[int8]("999999", stringx.Int8))
	fmt.Println(stringx.Str2Num[int16]("999999", stringx.Int16))
	fmt.Println(stringx.Str2Num[int32]("999999", stringx.Int32))
	fmt.Println(stringx.Str2Num[int64]("999999", stringx.Int64))
	fmt.Println(stringx.Str2Num[float32]("999999", stringx.Float32))
	fmt.Println(stringx.Str2Num[float64]("999999", stringx.Float64))
	fmt.Println(stringx.Str2Num[uint]("999999", stringx.Uint))
	fmt.Println(stringx.Str2Num[int]("999999", stringx.Int))
}

func TestReplace(t *testing.T) {
	fmt.Printf("%#v\n", stringx.Split("'A','B','大'", ",").Replace("'").Strings())
}
