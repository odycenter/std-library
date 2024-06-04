package internal_http

import (
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

type IPv4Ranges struct {
	ranges []uint32
}

func NewIPv4Ranges(cidrs []string) *IPv4Ranges {
	ranges := make([][]uint32, len(cidrs))
	for i, cidr := range cidrs {
		range_ := comparableIPRanges(cidr)
		ranges[i] = range_
	}
	return &IPv4Ranges{ranges: mergeRanges(ranges)}
}

func comparableIPRanges(cidr string) []uint32 {
	index := strings.Index(cidr, "/")
	if index <= 0 || index >= len(cidr)-1 {
		log.Panicf("invalid cidr, value=%s", cidr)
	}
	ip := net.ParseIP(cidr[:index])
	if ip == nil {
		log.Panicf("invalid ip address in cidr, value=%s", cidr)
	}
	address := toInteger(ip.To4())
	maskBits, err := strconv.Atoi(cidr[index+1:])
	if err != nil || maskBits < 0 || maskBits > 32 {
		log.Panicf("invalid mask bits in cidr, value=%s", cidr)
	}
	var mask uint32
	if maskBits == 0 {
		mask = 0
	} else {
		mask = ^uint32(0) << (32 - maskBits)
	}
	lowestIP := address & mask
	highestIP := address | ^mask
	return []uint32{lowestIP, highestIP}
}

func toInteger(ip net.IP) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func mergeRanges(ranges [][]uint32) []uint32 {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})
	results := make([]uint32, len(ranges)*2)
	index := 0
	for _, range_ := range ranges {
		if index > 1 && results[index-1] >= range_[0] {
			if results[index-1] < range_[1] {
				results[index-1] = range_[1]
			}
		} else {
			results[index] = range_[0]
			index++
			results[index] = range_[1]
			index++
		}
	}
	if index < len(results) {
		return results[:index]
	}
	return results
}

func withinRanges(ranges []uint32, value uint32) bool {
	for i := 0; i < len(ranges); i += 2 {
		if value >= ranges[i] && value <= ranges[i+1] {
			return true
		}
	}
	return false
}

func (r *IPv4Ranges) Matches(ip net.IP) bool {
	if len(r.ranges) == 0 {
		return false
	}
	comparableIP := toInteger(ip.To4())
	return withinRanges(r.ranges, comparableIP)
}
