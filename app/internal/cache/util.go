package internal_cache

import (
	reflects "std-library/reflect"
	"strings"
)

func cacheName(obj interface{}) string {
	return strings.ToLower(reflects.StructName(obj))
}
