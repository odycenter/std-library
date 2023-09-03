package maps

import (
	"sync"
)

// RWMutex 互斥锁实现的线程安全的map
type RWMutex[K T, V any] struct {
	m map[K]V
	l sync.RWMutex
}

// New 创建
func (m *RWMutex[K, V]) new() Inf[K, V] {
	return &RWMutex[K, V]{
		m: make(map[K]V),
	}
}

// Set 插入数据
func (m *RWMutex[K, V]) Set(k K, v V) {
	m.l.Lock()
	m.m[k] = v
	m.l.Unlock()
}

// Get 获取数据
func (m *RWMutex[K, V]) Get(k K) (V, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	v, ok := m.m[k]
	return v, ok
}

// Delete 根据k删除某条数据
func (m *RWMutex[K, V]) Delete(k K) {
	m.l.Lock()
	delete(m.m, k)
	m.l.Unlock()
}

// Len 获取map长度
func (m *RWMutex[K, V]) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return len(m.m)
}

// Range 循环map并执行操作
// Warning 不要在fn中执行较长耗时的任务，否则有可能造成进程阻塞
func (m *RWMutex[K, V]) Range(fn func(k K, v V) error) (err error) {
	m.l.RLock()
	defer m.l.RUnlock()
	for k, v := range m.m {
		err = fn(k, v)
		if err != nil {
			return err
		}
	}
	return err
}
