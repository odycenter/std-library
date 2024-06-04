package internal_web

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPHost struct {
	Host string
	Port int
}

func Parse(value string) *HTTPHost {
	index := strings.Index(value, ":")
	if index > 0 {
		port, _ := strconv.Atoi(value[index+1:])
		return &HTTPHost{
			Host: value[:index],
			Port: port,
		}
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return &HTTPHost{
		Host: "0.0.0.0",
		Port: port,
	}
}

func (h *HTTPHost) String() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}
