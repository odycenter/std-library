package slice

type T interface {
	any
}

type Inf[V T] interface {
	new() Inf[V]
	Push(v V)
	Index(i int) (V, bool)
	Delete(i int)
	Len() int
	Range(fn func(v V) error) (err error)
}

// New 创建指定的map
func New[V T](i Inf[V]) Inf[V] {
	return i.new()
}
