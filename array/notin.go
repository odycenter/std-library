package array

// Not 判断是否v不存在于s中
// 时间复杂度(O(1))
// s:[1,2,3,4]
// v:7 true
// v:2 false
func Not[T comparable](s []T, v T) bool {
	for _, t := range s {
		if t == v {
			return false
		}
	}
	return true
}

// NoOne 判断是否v中每个元素都不存在于s中
// 时间复杂度(O(1))
// s:[1,2,3,4]
// v1:[5,6,7] true
// v2:[4,5,6] false
func NoOne[T comparable](s []T, v []T) bool {
	var m = make(map[T]struct{})
	for _, e := range s {
		m[e] = struct{}{}
	}
	for _, e := range v {
		if _, ok := m[e]; ok {
			return false
		}
	}
	return true
}
