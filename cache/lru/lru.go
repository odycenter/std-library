package lru

import (
	"container/list"
	"sync"
)

type Cache[K comparable, V any] struct {
	cap  int
	keys map[K]*list.Element
	list *list.List
	sync.RWMutex
}

type entry[K comparable, V any] struct {
	K K
	V V
}

// New 创建LRU缓存
func New[K comparable, V any](cap int) *Cache[K, V] {
	return &Cache[K, V]{
		cap:  cap,
		keys: make(map[K]*list.Element),
		list: list.New(),
	}
}

// Put 插入
func (cache *Cache[K, V]) Put(key K, val V) {
	cache.Lock()
	defer cache.Unlock()
	if el, ok := cache.keys[key]; ok {
		el.Value = &entry[K, V]{
			K: key,
			V: val,
		}
		cache.list.MoveToFront(el)
	} else {
		el := cache.list.PushFront(&entry[K, V]{
			K: key,
			V: val,
		})
		cache.keys[key] = el
	}
	if cache.list.Len() > cache.cap {
		el := cache.list.Back()
		cache.list.Remove(el)
		delete(cache.keys, el.Value.(*entry[K, V]).K)
	}
}

// Get 获取，如果获取不到返回nil
func (cache *Cache[K, V]) Get(key K) (val V, ok bool) {
	cache.Lock()
	defer cache.Unlock()
	if el, ok := cache.keys[key]; ok {
		cache.list.MoveToFront(el)
		return el.Value.(*entry[K, V]).V, true
	}
	return
}

// Resize 重新設置大小
func (cache *Cache[K, V]) Resize(newSize int) {
	if newSize <= 0 {
		return
	}
	cache.Lock()
	defer cache.Unlock()
	if newSize < cache.list.Len() {
		for i := cache.list.Len() - newSize; i > 0; i-- {
			el := cache.list.Back()
			cache.list.Remove(el)
			delete(cache.keys, el.Value.(*entry[K, V]).K)
		}
	}
	cache.cap = newSize
}
