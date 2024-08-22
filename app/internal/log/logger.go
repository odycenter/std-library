package internal_log

import (
	"log"
	"os"
	"strings"
	"sync"
)

var (
	MaskedFields    []string
	maskedFieldsMap = make(map[string]struct{})
	maskedFieldsMu  sync.RWMutex
	Logger          = log.New(os.Stdout, "", 0)
)

func AddMaskedField(fieldNames ...string) {
	maskedFieldsMu.Lock()
	defer maskedFieldsMu.Unlock()

	for _, fieldName := range fieldNames {
		if _, exists := maskedFieldsMap[fieldName]; !exists {
			maskedFieldsMap[fieldName] = struct{}{}
			maskedFieldsMap[strings.ToLower(fieldName)] = struct{}{}
			MaskedFields = append(MaskedFields, fieldName)
		}
	}
}

func IsMaskedField(fieldName string) bool {
	_, ok := maskedFieldsMap[strings.ToLower(fieldName)]
	return ok
}
