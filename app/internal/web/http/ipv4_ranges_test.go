package internal_http

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithinRanges(t *testing.T) {
	ranges := []uint32{1, 1, 2, 3, 5, 8, 20, 30}

	assert.False(t, withinRanges(ranges, 0))
	assert.True(t, withinRanges(ranges, 2))
	assert.True(t, withinRanges(ranges, 3))
	assert.False(t, withinRanges(ranges, 4))
	assert.False(t, withinRanges(ranges, 9))
	assert.True(t, withinRanges(ranges, 22))
	assert.False(t, withinRanges(ranges, 31))
}

func TestMergeRanges(t *testing.T) {
	ranges1 := [][]uint32{{1, 2}, {4, 5}, {5, 9}, {3, 5}, {10, 20}}
	assert.Equal(t, []uint32{1, 2, 3, 9, 10, 20}, mergeRanges(ranges1))

	ranges2 := [][]uint32{{1, 2}, {10, 20}, {3, 5}, {30, 40}}
	assert.Equal(t, []uint32{1, 2, 3, 5, 10, 20, 30, 40}, mergeRanges(ranges2))
}

func TestMatchesAll(t *testing.T) {
	ranges := NewIPv4Ranges([]string{"0.0.0.0/0"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.True(t, ranges.Matches(net.ParseIP("127.0.0.1")))
	assert.True(t, ranges.Matches(net.ParseIP("10.10.0.1")))
}

func TestMatches(t *testing.T) {
	ranges := NewIPv4Ranges([]string{"192.168.1.0/24"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.254")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.2.1")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.0.1")))

	ranges = NewIPv4Ranges([]string{"192.168.1.1/32"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.1.2")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.1.3")))

	ranges = NewIPv4Ranges([]string{"192.168.1.1/31"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.0")))
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.1.2")))

	ranges = NewIPv4Ranges([]string{"192.168.1.1/30"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.0")))
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.2")))
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.3")))
	assert.False(t, ranges.Matches(net.ParseIP("192.168.1.4")))

	ranges = NewIPv4Ranges([]string{"119.137.52.0/22"})
	assert.True(t, ranges.Matches(net.ParseIP("119.137.52.1")))
	assert.True(t, ranges.Matches(net.ParseIP("119.137.53.1")))
	assert.True(t, ranges.Matches(net.ParseIP("119.137.53.254")))
	assert.True(t, ranges.Matches(net.ParseIP("119.137.54.254")))

	ranges = NewIPv4Ranges([]string{"42.200.0.0/16", "43.224.4.0/22", "43.224.28.0/22"})
	assert.False(t, ranges.Matches(net.ParseIP("42.119.0.1")))
	assert.True(t, ranges.Matches(net.ParseIP("42.200.218.1")))
	assert.False(t, ranges.Matches(net.ParseIP("42.201.218.1")))
	assert.False(t, ranges.Matches(net.ParseIP("43.224.32.1")))
}

func TestMatchWithEmptyRanges(t *testing.T) {
	ranges := NewIPv4Ranges([]string{})
	assert.False(t, ranges.Matches(net.ParseIP("192.168.1.1")))
}

func TestCIDRBoundaryValues(t *testing.T) {
	ranges := NewIPv4Ranges([]string{"0.0.0.0/32"})
	assert.False(t, ranges.Matches(net.ParseIP("0.0.0.1")))
	assert.True(t, ranges.Matches(net.ParseIP("0.0.0.0")))

	ranges = NewIPv4Ranges([]string{"255.255.255.255/0"})
	assert.True(t, ranges.Matches(net.ParseIP("192.168.1.1")))
	assert.True(t, ranges.Matches(net.ParseIP("127.0.0.1")))
	assert.True(t, ranges.Matches(net.ParseIP("10.10.0.1")))

	ctl := &IPv4AccessControl{}
	result := ctl.Validate("61.220.68.81")
	println(result)
}
