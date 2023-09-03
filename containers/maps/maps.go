package maps

type T interface {
	comparable
}

type Inf[K T, V any] interface {
	new() Inf[K, V]
	Set(k K, v V)
	Get(k K) (V, bool)
	Delete(k K)
	Len() int
	Range(fn func(k K, v V) error) (err error)
}

// New 创建指定的map
func New[K T, V any](i Inf[K, V]) Inf[K, V] {
	return i.new()
}
