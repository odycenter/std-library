// Package pagex 分页
package pagex

// Pagination 分页结构体
type Pagination[T Num] struct {
	CurrPage     T //当前页数
	PageSize     T //每页条数
	MaxCount     T //总行数
	MaxPageCount T //最大每页数量
}

type Num interface {
	int | int32 | int64
}

// New 创建分页
func New[T Num](currPage T, pageSize T, maxCount T) (res *Pagination[T]) {
	res = &Pagination[T]{}
	res.set(currPage, pageSize, maxCount)
	return
}

func (p *Pagination[T]) set(currPage T, pageSize T, maxCount T) {
	p.CurrPage = currPage
	p.PageSize = pageSize
	p.MaxCount = maxCount

	if p.PageSize <= 0 {
		p.PageSize = 15
	}
	if p.MaxCount <= 0 {
		p.MaxCount = 0
	}
	if p.CurrPage <= 0 {
		p.CurrPage = 1
	}
	p.MaxPageCount = p.MaxCount / p.PageSize
	if p.MaxCount%p.PageSize > 0 {
		p.MaxPageCount += 1
	}
	if p.MaxPageCount <= 0 {
		p.MaxPageCount = 1
	}
}

// Offset 获取偏移值
func (p *Pagination[T]) Offset() (offset T) {
	if p.CurrPage > p.MaxPageCount {
		return p.MaxCount
	}
	return (p.CurrPage - 1) * p.PageSize
}
