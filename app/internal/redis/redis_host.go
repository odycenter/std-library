package internal_redis

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var DefaultPort = 6379

type RedisHost struct {
	host string
	port int
}

func Host(host string) *RedisHost {
	index := strings.Index(host, ":")
	if index == 0 || index == len(host)-1 {
		log.Panic("invalid host format, host=" + host)
	}
	h := RedisHost{}
	if index != -1 {
		h.host = host[:index]
		var err error
		h.port, err = strconv.Atoi(host[index+1:])
		if err != nil {
			log.Panic("invalid host format, host=" + host)
		}
	} else {
		h.host = host
		h.port = DefaultPort
	}

	return &h
}

func (h *RedisHost) String() string {
	return h.host + ":" + fmt.Sprintf("%d", h.port)
}
