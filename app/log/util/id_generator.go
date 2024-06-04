package util

import (
	"encoding/hex"
	"hash/crc32"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// lowOrderThreeBytes is a mask for the low order three bytes of an int32.
	lowOrderThreeBytes = 0x00ffffff
)

var once sync.Once
var instance IDGenerator

func generateMachineID() (machineID int, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	var str string
	for _, iface := range interfaces {
		str = str + iface.Name

		mac := iface.HardwareAddr
		if len(mac) > 0 {
			str = str + hex.EncodeToString(mac)
		}
	}

	str = str + strconv.Itoa(rand.Int())
	return String(str), nil
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// a non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

type IDGenerator struct {
	counter           atomic.Int32
	machineIdentifier int32
}

func GetIDGenerator() *IDGenerator {
	once.Do(func() {
		machineID, err := generateMachineID()
		if err != nil {
			return
		}
		instance = IDGenerator{
			machineIdentifier: int32(machineID) & lowOrderThreeBytes,
		}
		instance.counter.Store(rand.Int31())
	})
	return &instance
}

func (g *IDGenerator) Next(now time.Time) string {
	milli := now.UnixMilli()
	counter := g.counter.Add(1) & lowOrderThreeBytes
	machineIdentifier := g.machineIdentifier

	buf := make([]byte, 10)
	buf[0] = byte(milli >> 32) // save 5 bytes time in ms, about 34 years value space
	buf[1] = byte(milli >> 24)
	buf[2] = byte(milli >> 16)
	buf[3] = byte(milli >> 8)
	buf[4] = byte(milli)
	buf[5] = byte(machineIdentifier >> 16) // 3 bytes as machine id, about 16M value space
	buf[6] = byte(machineIdentifier >> 8)
	buf[7] = byte(machineIdentifier)
	buf[8] = byte(counter >> 8) // 2 bytes for max 65k actions per ms per server
	buf[9] = byte(counter)
	return hex.EncodeToString(buf)
}
