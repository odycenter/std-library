package array

// In 判断v是否存在于s中
// 时间复杂度(O(n))
func In[T comparable](s []T, v T) bool {
	for _, t := range s {
		if t == v {
			return true
		}
	}
	return false
}

// Index 查找v在s中第一次出现的位置
// 时间复杂度(O(n))
func Index[T comparable](s []T, v T) int {
	for i, t := range s {
		if t == v {
			return i
		}
	}
	return -1
}

// All 查找v在s中所有出现的所有位置，如果存在则返回存在位置的所有index的数组,不存在则返回空数组
// 时间复杂度(O(n))
func All[T comparable](s []T, v T) []int {
	var a []int
	for i, t := range s {
		if t == v {
			a = append(a, i)
		}
	}
	return a
}

// Last 查找v在s中所有出现的最后位置，如果不存在则返回-1
// 时间复杂度(O(n))
func Last[T comparable](s []T, v T) int {
	var last = -1
	for i, t := range s {
		if t == v {
			last = i
		}
	}
	return last
}
