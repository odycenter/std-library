package array

import (
	"fmt"
	"log"
	"strconv"
)

// Nums2Strings 数字Slice转字符串Slice
func Nums2Strings[T Numeric](s []T) []string {
	var r []string
	for _, num := range s {
		r = append(r, fmt.Sprint(num))
	}
	return r
}

// Strings2Nums 字符串Slice转数字Slice
func Strings2Nums[T Numeric](s []string, t convertTo) []T {
	switch t {
	case Uint8:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(uint8(i))
		})
	case Uint16:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(uint16(i))
		})
	case Uint32:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(uint32(i))
		})
	case Uint64:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(uint64(i))
		})
	case Int8:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(int8(i))
		})
	case Int16:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(int16(i))
		})
	case Int32:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(int32(i))
		})
	case Int64:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(i)
		})
	case Float32:
		return each(s, func(e string) T {
			f, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(float32(f))
		})
	case Float64:
		return each(s, func(e string) T {
			f, err := strconv.ParseFloat(e, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(f)
		})
	case Int:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(int(i))
		})
	case Uint:
		return each(s, func(e string) T {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Panicln(err.Error())
			}
			return T(uint(i))
		})
	}
	return nil
}

// Bytes2String Bytes2String 将byte数组，作为普通数组转换为字符串，而非按 ASCII 码转换
// []byte{97,98} => "9798"
func Bytes2String(bs []byte) (s string) {
	for _, b := range bs {
		s += fmt.Sprint(b)
	}
	return
}

// Map Slice转map
func Map[T Numeric | string](s []T) map[T]struct{} {
	var m = make(map[T]struct{})
	for _, e := range s {
		m[e] = struct{}{}
	}
	return m
}

func each[T Numeric](s []string, fn func(e string) T) (rs []T) {
	for _, e := range s {
		rs = append(rs, fn(e))
	}
	return rs
}

// Flatten 将map展开为对应类型数组
// map[string]int => []int
func Flatten[K comparable, V any](m map[K]V) (vs []V) {
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// FlattenAny 将map展开为interface数组
// map[string]int => []interface{}
func FlattenAny[K comparable, V any](m map[K]V) (vs []any) {
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}
