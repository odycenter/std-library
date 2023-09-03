package slice

import "sync"

type RWMutex[V T] struct {
	s []V
	l sync.RWMutex
}

func (s *RWMutex[V]) new() Inf[V] {
	return &RWMutex[V]{
		s: []V{},
	}
}

// Push 插入元素到最后
func (s *RWMutex[V]) Push(v V) {
	s.l.Lock()
	s.s = append(s.s, v)
	s.l.Unlock()
}

// Index 按Index查找元素
func (s *RWMutex[V]) Index(i int) (V, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.s[i], true
}

// Delete 删除第i个元素
func (s *RWMutex[V]) Delete(i int) {
	s.l.Lock()
	l := len(s.s)
	if i >= l || i < 0 {
		return
	}
	s.s = append(s.s[:i], s.s[i+1:]...)
	s.l.Unlock()
}

// Len 元素个数
func (s *RWMutex[V]) Len() int {
	s.l.RLock()
	l := len(s.s)
	s.l.RUnlock()
	return l
}

// Range 循环slice并执行操作
// Warning 不要在fn中执行较长耗时的任务，否则有可能造成进程阻塞
func (s *RWMutex[V]) Range(fn func(v V) error) (err error) {
	s.l.RLock()
	defer s.l.RUnlock()
	for _, v := range s.s {
		err = fn(v)
		if err != nil {
			return err
		}
	}
	return err
}
