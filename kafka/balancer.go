package kafka

import (
	"github.com/odycenter/std-library/app/log/util"
	"github.com/segmentio/kafka-go"
	"sync"
)

type KeyHashBalancer struct {
	rr     RoundRobinBalancer
	offset uint32
}

func (pf *KeyHashBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	return pf.BalanceByKey(msg.Key, partitions...)
}

func (pf *KeyHashBalancer) BalanceByKey(keyByteArray []byte, partitions ...int) (partition int) {
	if keyByteArray == nil || len(keyByteArray) == 0 {
		return pf.rr.Balance(partitions...)
	}

	key := string(keyByteArray)
	hashcode := util.String(key)
	keyHash := uint32(hashcode)
	return pf.balance(keyHash, partitions)
}

func (pf *KeyHashBalancer) balance(keyHash uint32, partitions []int) int {
	length := uint32(len(partitions))
	return partitions[keyHash%length]
}

type RoundRobinBalancer struct {
	ChunkSize int
	// Use a 32 bits integer so RoundRobin values don't need to be aligned to
	// apply increments.
	counter uint32

	mutex sync.Mutex
}

// Balance satisfies the Balancer interface.
func (rr *RoundRobinBalancer) Balance(partitions ...int) int {
	return rr.balance(partitions)
}

func (rr *RoundRobinBalancer) balance(partitions []int) int {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	if rr.ChunkSize < 1 {
		rr.ChunkSize = 1
	}

	length := len(partitions)
	counterNow := rr.counter
	offset := int(counterNow / uint32(rr.ChunkSize))
	rr.counter++
	return partitions[offset%length]
}
