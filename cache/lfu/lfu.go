package lfu

import "sync"

type Node[K comparable, V any] struct {
	key        K
	value      V
	prev, next *Node[K, V]
	freqNode   *freqNode[K, V]
}

type List[K comparable, V any] struct {
	head, tail *Node[K, V]
	total      int
}

type freqNode[K comparable, V any] struct {
	freq       int
	dl         *List[K, V]
	prev, next *freqNode[K, V]
}

type freqList[K comparable, V any] struct {
	head *freqNode[K, V]
	tail *freqNode[K, V]
}

func (fl *freqList[K, V]) removeNode(node *freqNode[K, V]) {
	node.next.prev = node.prev
	node.prev.next = node.next
}

func (fl *freqList[K, V]) lastFreq() *freqNode[K, V] {
	return fl.tail.prev
}

func (fl *freqList[K, V]) addNode(node *Node[K, V]) {
	if fqNode := fl.lastFreq(); fqNode.freq == 1 {
		node.freqNode = fqNode
		fqNode.dl.addToHead(node)
	} else {
		newNode := &freqNode[K, V]{
			freq: 1,
			dl:   initList[K, V](),
		}

		node.freqNode = newNode
		newNode.dl.addToHead(node)

		fqNode.next = newNode
		newNode.prev = fqNode
		newNode.next = fl.tail
		fl.tail.prev = newNode
	}
}

func (l *List[K, V]) isEmpty() bool {
	return l.total == 0
}

func (l *List[K, V]) GetTotal() int {
	return l.total
}

func (l *List[K, V]) addToHead(node *Node[K, V]) {
	node.next = l.head.next
	node.prev = l.head
	l.head.next.prev = node
	l.head.next = node
	l.total++
}

func (l *List[K, V]) removeNode(node *Node[K, V]) {
	node.next.prev = node.prev
	node.prev.next = node.next
	l.total--
}

func (l *List[K, V]) moveToHead(node *Node[K, V]) {
	l.removeNode(node)
	l.addToHead(node)
}

func (l *List[K, V]) removeTail() *Node[K, V] {
	node := l.tail.prev
	l.removeNode(node)
	return node
}

func initNode[K comparable, V any](k K, v V) *Node[K, V] {
	return &Node[K, V]{
		key:   k,
		value: v,
	}
}

func initList[K comparable, V any]() *List[K, V] {
	l := List[K, V]{
		head: new(Node[K, V]),
		tail: new(Node[K, V]),
	}
	l.head.next = l.tail
	l.tail.prev = l.head

	return &l
}

type Cache[K comparable, V any] struct {
	cache     map[any]*Node[K, V]
	size, cap int
	freqList  *freqList[K, V]
	sync.RWMutex
}

// New 创建LFU缓存
func New[K comparable, V any](cap int) *Cache[K, V] {
	ca := &Cache[K, V]{
		cap:   cap,
		cache: make(map[any]*Node[K, V]),
	}
	ca.freqList = &freqList[K, V]{
		head: &freqNode[K, V]{},
		tail: &freqNode[K, V]{},
	}
	ca.freqList.head.next = ca.freqList.tail
	ca.freqList.tail.prev = ca.freqList.head
	return ca
}

func (cache *Cache[K, V]) incrFreq(node *Node[K, V]) {
	curFreqNode := node.freqNode
	curNode := curFreqNode.dl

	if curFreqNode.prev.freq == curFreqNode.freq+1 {
		curNode.removeNode(node)
		curFreqNode.prev.dl.addToHead(node)
		node.freqNode = curFreqNode.prev
	} else if curNode.GetTotal() == 1 {
		curFreqNode.freq++
	} else {
		curNode.removeNode(node)
		newFreqNode := &freqNode[K, V]{
			freq: curFreqNode.freq + 1,
			dl:   initList[K, V](),
		}
		newFreqNode.dl.addToHead(node)
		node.freqNode = newFreqNode
		newFreqNode.next = curFreqNode
		newFreqNode.prev = curFreqNode.prev
		curFreqNode.prev.next = newFreqNode
		curFreqNode.prev = newFreqNode
	}

	if curNode.isEmpty() {
		cache.freqList.removeNode(curFreqNode)
	}
}

// Put 插入
func (cache *Cache[K, V]) Put(key K, value V) {
	cache.Lock()
	defer cache.Unlock()
	if cache.cap == 0 {
		return
	}
	if n, ok := cache.cache[key]; ok {
		n.value = value
		cache.incrFreq(n)
	} else {
		//因为有根节点元素，所以长度多一个
		if cache.size+1 >= cache.cap {
			fqNode := cache.freqList.lastFreq()
			node := fqNode.dl.removeTail()
			cache.size--
			delete(cache.cache, node.key)
		}

		newNode := initNode[K, V](key, value)
		cache.cache[key] = newNode
		cache.freqList.addNode(newNode)
		cache.size++
	}
}

// Get 获取，如果获取不到返回nil
func (cache *Cache[K, V]) Get(key K) (val V, ok bool) {
	cache.RLock()
	defer cache.RUnlock()
	if n, ok := cache.cache[key]; ok {
		cache.incrFreq(n)
		return n.value, true
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
	if newSize < cache.size {
		for i := cache.size - newSize; i > 0; i-- {
			fqNode := cache.freqList.lastFreq()
			//if fqNode.next.dl == nil {
			//	break
			//}
			node := fqNode.dl.removeTail()
			cache.size--
			delete(cache.cache, node.key)
		}
	}
	cache.cap = newSize
}
