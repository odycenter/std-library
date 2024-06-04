package property

import (
	"strings"
)

func EnvVarName(str string) string {
	var result []int32
	for _, r := range str {
		if r == '.' {
			result = append(result, '_')
		} else {
			result = append(result, r)
		}
	}
	return strings.ToUpper(string(result))
}
