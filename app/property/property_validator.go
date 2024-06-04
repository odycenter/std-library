package property

import "log"

type Validator struct {
	usedProperties Properties
}

func NewValidator() *Validator {
	return &Validator{usedProperties: make(Properties)}
}

func (p *Validator) Validate(keys []string) {
	notUsedKeys := make([]string, 0)
	for _, key := range keys {
		_, ok := p.usedProperties[key]
		if !ok {
			notUsedKeys = append(notUsedKeys, key)
		}
	}
	if len(notUsedKeys) > 0 {
		log.Panic("not used properties: ", notUsedKeys)
	}
}

func (p *Validator) Add(key string) {
	p.usedProperties[key] = key
}
