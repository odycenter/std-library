package module

import (
	"strings"
)

type IPv4RangePropertyValueParser struct {
	value string
}

func NewIPv4RangePropertyValueParser(value string) *IPv4RangePropertyValueParser {
	return &IPv4RangePropertyValueParser{value: value}
}

// Parse supports two formats:
// 1. cidr,cidr
// 2. name1: cidr, cidr; name2: cidr
func (p *IPv4RangePropertyValueParser) Parse() []string {
	if strings.TrimSpace(p.value) == "" {
		return []string{}
	}
	if strings.Index(p.value, ":") > 0 {
		return p.parseSemicolonDelimited(p.value)
	}
	return p.parseCommaDelimited(p.value)
}

func (p *IPv4RangePropertyValueParser) parseCommaDelimited(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}

func (p *IPv4RangePropertyValueParser) parseSemicolonDelimited(value string) []string {
	var results []string
	items := strings.Split(value, ";")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		index := strings.Index(item, ":")
		if index == -1 {
			continue
		}
		parts := strings.Split(item[index+1:], ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				results = append(results, part)
			}
		}
	}
	return results
}
