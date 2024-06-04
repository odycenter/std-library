package kafka

import "strings"

type Uri struct {
	uri string
}

func (u *Uri) Uri(uri string) {
	u.uri = uri
}

func (u *Uri) Parse() []string {
	values := strings.Split(u.uri, ",")
	uris := make([]string, len(values))
	for i, value := range values {
		value = strings.TrimSpace(value)
		if !strings.Contains(value, ":") {
			value += ":9092"
		}
		uris[i] = value
	}
	return uris
}
