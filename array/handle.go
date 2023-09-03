package array

import "strings"

// Replace 替换字符串中的匹配字符
// default r ""
// Warning 非线程安全
func Replace(v []string, s string, r ...string) []string {
	if len(r) == 0 {
		r = append(r, "")
	}
	for i, e := range v {
		v[i] = strings.Replace(e, s, r[0], -1)
	}
	return v
}

// Reverse 将当前数组反向
func Reverse[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// RemoveDuplicates 数组去重
func RemoveDuplicates[T comparable](slice []T) []T {
	mUnique := make(map[T]bool)
	var unique []T
	for _, num := range slice {
		if _, ok := mUnique[num]; !ok {
			unique = append(unique, num)
			mUnique[num] = true
		}
	}
	return unique
}

// Remove 删除数组和elem相同的元素
// n指定删除的数量，-1删除全部
func Remove[T comparable](slice []T, elem T, n int) []T {
	var ret []T
	for _, t := range slice {
		if (n == -1 || n > 0) && t == elem {
			if n != -1 {
				n--
			}
			continue
		}
		ret = append(ret, t)
	}
	return ret
}

// RemoveFn 删除数组中传入fn返回值为true的元素
// n指定删除的数量，-1删除全部
func RemoveFn[T any](slice []T, fn func(elem T) bool) []T {
	var ret []T
	for _, t := range slice {
		if fn(t) {
			continue
		}
		ret = append(ret, t)
	}
	return ret
}
