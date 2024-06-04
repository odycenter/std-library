package util

import "strings"

func Filter(value string, maskedFields ...string) string {
	result := value
	for _, maskedField := range maskedFields {
		current := -1
		for {
			current = indexOf(result, "\""+maskedField+"\"", current)
			if current < 0 {
				break
			}
			rangeArray := maskRange(&result, current+len(maskedField)+2)
			if rangeArray == nil {
				break
			}
			result = replace(result, rangeArray[0], rangeArray[1])
			current = rangeArray[1]
		}
	}
	return result
}

func replace(input string, start int, length int) string {
	return strings.Join([]string{input[:start], "******", input[length:]}, "")
}

func indexOf(str string, substring string, fromIndex int) int {
	if fromIndex < 0 {
		fromIndex = 0
	}
	if fromIndex >= len(str) {
		return -1
	}
	index := strings.Index(str[fromIndex:], substring)
	if index < 0 {
		return index
	}
	return index + fromIndex
}

func maskRange(builder *string, start int) []int {
	escaped := false
	maskStart := -1
	length := len(*builder)
	for index := start; index < length; index++ {
		ch := (*builder)[index]
		if ch == '\\' {
			escaped = true
		} else if !escaped && maskStart < 0 && ch == '"' {
			maskStart = index + 1
		} else if !escaped && maskStart >= 0 && ch == '"' {
			return []int{maskStart, index}
		} else if maskStart < 0 && (ch == ',' || ch == '{' || ch == '}' || ch == '[' || ch == ']') {
			return nil // not found start double quote, and reached next field
		} else {
			escaped = false
		}

	}
	if maskStart >= 0 {
		return []int{maskStart, length}
	}
	return nil
}

func ShouldMask(value string, maskedFields ...string) bool {
	if !strings.HasPrefix(value, "{") {
		return false
	}
	for _, maskedField := range maskedFields {
		index := strings.Index(value, "\""+maskedField+"\"")
		if index >= 0 {
			return true
		}
	}
	return false
}
