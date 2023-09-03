package maps

import (
	"fmt"
	"sync"
)

// Shards 分段map
// 为提升对大量数据
type Shards[K T, V any] []*Shard[K, V]

type Shard[K T, V any] struct {
	m map[K]V //container
	l sync.RWMutex
}

func (m *Shards[K, V]) new() Inf[K, V] {
	l := make(Shards[K, V], shardCount)
	for i := 0; i < shardCount; i++ {
		l[i] = &Shard[K, V]{}
		l[i].m = make(map[K]V)
	}
	return &l
}

const (
	offset32   = uint32(2166136261)
	prime32    = 16777619
	shardCount = 32
)

// fnv32Hash算法
func fnv32[K T](k K) uint32 {
	sk := fmt.Sprint(k)
	h := offset32
	for i := 0; i < len(sk); i++ {
		h *= prime32
		h ^= uint32(sk[i])
	}
	return h
}

// 获取当前key所在的分片
func (m *Shards[K, V]) getShard(k K) *Shard[K, V] {
	return (*m)[uint(fnv32(k))%uint(shardCount)]
}

func (m *Shards[K, V]) Get(k K) (V, bool) {
	s := m.getShard(k)
	s.l.RLock()
	v, ok := s.m[k]
	s.l.RUnlock()
	return v, ok
}

func (m *Shards[K, V]) Set(k K, v V) {
	s := m.getShard(k)
	s.l.Lock()
	s.m[k] = v
	s.l.Unlock()
}

func (m *Shards[K, V]) Delete(k K) {
	s := m.getShard(k)
	s.l.RLock()
	delete(s.m, k)
	s.l.RUnlock()
}

func (m *Shards[K, V]) Len() int {
	l := 0
	wg := sync.WaitGroup{}
	wg.Add(shardCount)
	for i := 0; i < shardCount; i++ {
		go func(index int) {
			s := (*m)[index]
			s.l.RLock()
			l += len(s.m)
			s.l.RUnlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	return l
}

// Range 循环map并执行操作
// Warning 不要在fn中执行较长耗时的任务，否则有可能造成进程阻塞
func (m *Shards[K, V]) Range(fn func(k K, v V) error) (err error) {
	for _, s := range *m {
		s.l.RLock()
		for k, v := range s.m {
			err := fn(k, v)
			if err != nil {
				return err
			}
		}
		s.l.RUnlock()
	}
	return err
}
