// Package stringx 字符串操作扩展
package stringx

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"strconv"
	"strings"
)

// Str2Num 字符串转数字
func Str2Num[T Numeric](s string, t convertTo) T {
	if s == "" {
		return 0
	}
	switch t {
	case Uint8:
		return T(uint8(toInt(s)))
	case Uint16:
		return T(uint16(toInt(s)))
	case Uint32:
		return T(uint32(toInt(s)))
	case Uint64:
		return T(uint64(toInt(s)))
	case Int8:
		return T(int8(toInt(s)))
	case Int16:
		return T(int16(toInt(s)))
	case Int32:
		return T(int32(toInt(s)))
	case Int64:
		return T(toInt(s))
	case Float32:
		return T(float32(toFloat(s)))
	case Float64:
		return T(toFloat(s))
	case Int:
		return T(int(toInt(s)))
	case Uint:
		return T(uint(toInt(s)))
	}
	return 0
}

// Str2NumR 字符串转数字并返回是否成功
func Str2NumR[T Numeric](s string, t convertTo) (T, bool) {
	if !isNumeric(s) {
		return 0, false
	}
	return Str2Num[T](s, t), true
}

func isNumeric(s string) bool {
	var (
		dotCount = 0
		length   = len(s)
	)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if s[i] == '-' && i == 0 {
			continue
		}
		if s[i] == '.' {
			dotCount++
			if i > 0 && i < length-1 {
				continue
			} else {
				return false
			}
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return dotCount <= 1
}

func toInt(e string) int64 {
	i, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		log.Panicln(err.Error())
	}
	return i
}
func toFloat(e string) float64 {
	f, err := strconv.ParseFloat(e, 64)
	if err != nil {
		log.Panicln(err.Error())
	}
	return f
}

type T struct {
	v []string
}

// Strings 输出-字符串数组
func (t *T) Strings() []string {
	return t.v
}

// Int32s 输出-int32数组
func (t *T) Int32s() []int32 {
	if len(t.v) == 0 {
		return []int32{}
	}
	var elems []int32
	for _, v := range t.v {
		elem, err := strconv.ParseInt(v, 0, 0)
		if err != nil {
			return []int32{}
		}
		elems = append(elems, int32(elem))
	}
	return elems
}

// Int64s 输出-int64数组
func (t *T) Int64s() []int64 {
	if len(t.v) == 0 {
		return []int64{}
	}
	elems := make([]int64, 0, len(t.v))
	for _, v := range t.v {
		elem, err := strconv.ParseInt(v, 0, 0)
		if err != nil {
			return []int64{}
		}
		elems = append(elems, elem)
	}
	return elems
}

// Ints 输出-int数组
func (t *T) Ints() []int {
	if len(t.v) == 0 {
		return []int{}
	}
	var elems []int
	for _, v := range t.v {
		elem, err := strconv.Atoi(v)
		if err != nil {
			return []int{}
		}
		elems = append(elems, elem)
	}
	return elems
}

// Int8s 输出-int8数组
func (t *T) Int8s() []int8 {
	if len(t.v) == 0 {
		return []int8{}
	}
	var elems []int8
	for _, v := range t.v {
		elem, err := strconv.Atoi(v)
		if err != nil {
			return []int8{}
		}
		elems = append(elems, int8(elem))
	}
	return elems
}

// Float64s 输出-float64数组
func (t *T) Float64s() []float64 {
	if len(t.v) == 0 {
		return []float64{}
	}
	var elems []float64
	for _, v := range t.v {
		elem, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return []float64{}
		}
		elems = append(elems, elem)
	}
	return elems
}

// Float32s 输出-float32数组
func (t *T) Float32s() []float32 {
	if len(t.v) == 0 {
		return []float32{}
	}
	var elems []float32
	for _, v := range t.v {
		elem, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return []float32{}
		}
		elems = append(elems, float32(elem))
	}
	return elems
}

// Trim 去除字符串中的首尾匹配字符
// Warning 非线程安全
func (t *T) Trim(s string) *T {
	for i, v := range t.v {
		t.v[i] = strings.Trim(v, s)
	}
	return t
}

// Replace 替换字符串中的匹配字符
// default r ""
// Warning 非线程安全
func (t *T) Replace(s string, r ...string) *T {
	if len(r) == 0 {
		r = append(r, "")
	}
	for i, v := range t.v {
		t.v[i] = strings.Replace(v, s, r[0], -1)
	}
	return t
}

// Num2Str 数字转换为string
func Num2Str[T Numeric](in T) string {
	return fmt.Sprint(in)
}

// Encode 字符串编码转换
func Encode(in []byte, charset charset) string {
	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(in)
		str = string(decodeBytes)
	case HZGB2312:
		decodeBytes, _ := simplifiedchinese.HZGB2312.NewDecoder().Bytes(in)
		str = string(decodeBytes)
	case GBK:
		decodeBytes, _ := simplifiedchinese.GBK.NewDecoder().Bytes(in)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(in)
	}
	return str
}

// Str2Bool 字符串转bool
func Str2Bool(in string) bool {
	b, _ := strconv.ParseBool(in)
	return b
}
