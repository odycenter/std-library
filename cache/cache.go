package cache

type Inf[K comparable, V any] interface {
	Put(K, V)
	Get(K) (any, bool)
	Resize(int)
}
