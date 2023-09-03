// Package rand 随机值生成器
package rand

import (
	"math/rand"
	"strings"
	"sync"
)

// R 结果类型适配结构
type R struct {
	*rand.Rand
	sync.Mutex
}

func (r *R) intN(n int) int {
	if r.Rand == nil {
		return rand.Intn(n)
	}
	r.Lock()
	defer r.Unlock()
	return r.Intn(n)
}

var (
	dirLetter  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	dirNum     = "0123456789"
	dirSpacial = "~!@#$%^&*()_+-=*.,;:"
)

// Rand 指定随机数生成器
// 使用默认 seed 时，可随意使用
// 使用传入 seed 时，请自行保存返回的 rand 对象,切勿每次调用 rand.Rand(val)
// r := rand.Rand(val)
// r1 := r.Range(1,10)
// r2 := r.Range(1,10)
func Rand(seed ...int64) *R {
	if len(seed) == 0 {
		return &R{}
	}
	return &R{Rand: rand.New(rand.NewSource(seed[0]))}
}

// Strings 字母和符号
func (r *R) Strings(l int) string {
	bytes := make([]byte, l)
	b := strings.Builder{}
	b.WriteString(dirLetter)
	b.WriteString(dirSpacial)
	dir := b.String()
	dl := len(dir)
	for i := 0; i < l; i++ {
		bytes[i] = dir[r.intN(dl)]
	}
	return string(bytes)
}

// Letters 仅字母
func (r *R) Letters(l int) string {
	bytes := make([]byte, l)
	dl := len(dirLetter)
	for i := 0; i < l; i++ {
		bytes[i] = dirLetter[r.intN(dl)]
	}
	return string(bytes)
}

// General 字母数字
func (r *R) General(l int) string {
	bytes := make([]byte, l)
	b := strings.Builder{}
	b.WriteString(dirLetter)
	b.WriteString(dirNum)
	dir := b.String()
	dl := len(dir)
	for i := 0; i < l; i++ {
		bytes[i] = dir[r.intN(dl)]
	}
	return string(bytes)
}

// Number 仅数字
func (r *R) Number(l int) string {
	bytes := make([]byte, l)
	b := strings.Builder{}
	b.WriteString(dirNum)
	dir := b.String()
	dl := len(dir)
	for i := 0; i < l; i++ {
		bytes[i] = dir[r.intN(dl)]
	}
	return string(bytes)
}

// Mixes 字母数字符号
func (r *R) Mixes(l int, dirEx string) string {
	bytes := make([]byte, l)
	b := strings.Builder{}
	b.WriteString(dirLetter)
	b.WriteString(dirSpacial)
	b.WriteString(dirNum)
	b.WriteString(dirEx)
	dir := b.String()
	dl := len(dir)
	for i := 0; i < l; i++ {
		bytes[i] = dir[r.intN(dl)]
	}
	return string(bytes)
}

// Custom 自定义
func (r *R) Custom(l int, dir string) string {
	bytes := make([]byte, l)
	dl := len(dir)
	for i := 0; i < l; i++ {
		bytes[i] = dir[r.intN(dl)]
	}
	return string(bytes)
}

// Range 范围
func (r *R) Range(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return r.intN(max-min) + min
}
