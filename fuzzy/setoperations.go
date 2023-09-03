package fuzzy

type StringSet struct {
	elements map[string]bool
}

func NewStringSet(slice []string) *StringSet {
	sliceStringSet := make(map[string]bool)
	for _, b := range slice {
		sliceStringSet[b] = true
	}
	s := new(StringSet)
	s.elements = sliceStringSet
	return s
}

// Difference 返回此集合中存在但其他集合中不存在的字符串集合
func (s *StringSet) Difference(other *StringSet) *StringSet {
	diff := new(StringSet)
	diff.elements = make(map[string]bool)
	for k, v := range s.elements {
		if _, ok := other.elements[k]; !ok {
			diff.elements[k] = v
		}
	}
	return diff
}

// Intersect 返回两个集合中都包含的字符串集合
func (s *StringSet) Intersect(other *StringSet) *StringSet {
	intersection := new(StringSet)
	intersection.elements = make(map[string]bool)
	for k, v := range s.elements {
		if _, ok := other.elements[k]; ok {
			intersection.elements[k] = v
		}
	}
	return intersection
}

// Equals 如果两个集合包含相同的元素，则返回 true
func (s *StringSet) Equals(other *StringSet) bool {
	if len(s.elements) != len(other.elements) {
		return false
	}

	for k, _ := range s.elements {
		if _, ok := other.elements[k]; !ok {
			return false
		}
	}
	return true
}

// ToSlice 从集合中生成一个字符串切片
func (s *StringSet) ToSlice() []string {
	keys := make([]string, len(s.elements))

	i := 0
	for k := range s.elements {
		keys[i] = k
		i++
	}
	return keys
}
