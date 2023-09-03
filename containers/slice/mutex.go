package slice

import "sync"

type Mutex[V T] struct {
	s []V
	l sync.Mutex
}

func (s *Mutex[V]) new() Inf[V] {
	return &Mutex[V]{
		s: []V{},
	}
}

// Push 插入元素到最后
func (s *Mutex[V]) Push(v V) {
	s.l.Lock()
	s.s = append(s.s, v)
	s.l.Unlock()
}

// Index 按Index查找元素
func (s *Mutex[V]) Index(i int) (V, bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.s[i], true
}

// Delete 删除第i个元素
func (s *Mutex[V]) Delete(i int) {
	s.l.Lock()
	l := len(s.s)
	if i >= l || i < 0 {
		s.l.Unlock()
		return
	}
	s.s = append(s.s[:i], s.s[i+1:]...)
	s.l.Unlock()
}

// Len 元素个数
func (s *Mutex[V]) Len() int {
	s.l.Lock()
	l := len(s.s)
	s.l.Unlock()
	return l
}

// Range 循环slice并执行操作
// Warning 不要在fn中执行较长耗时的任务，否则有可能造成进程阻塞
func (s *Mutex[V]) Range(fn func(v V) error) (err error) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, v := range s.s {
		err = fn(v)
		if err != nil {
			return err
		}
	}
	return err
}
