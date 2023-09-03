package elastic

func Int(v int) *int {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

// 查询条件类型
type queryKey int

// SearchQuery 查询条件集合。适用于高级查询
// Index 索引
// From 第几条记录
// Size 分页每页长度
// Sort 排序 {"score":true, "create_time":false} bool 是否正序
// Pretty 是否美化json
// QueryEntity ES查询参数
type SearchQuery struct {
	Index       []string
	From        *int
	Size        *int
	Pretty      *bool
	Sort        []Sort
	QueryEntity []QueryEntity
}

func (s *SearchQuery) getSize() *int {
	size := 500
	if s.Size == nil {
		return &size
	}
	return s.Size
}

// SimpleSearchQuery 查询条件集合。适用于简单查询
// Index 索引
// IDs LogIDs
// From 第几条记录
// Size 分页每页长度
// Sort 排序 {"score":true, "create_time":false} bool 是否正序
// Pretty 是否美化json
// QueryEntity ES查询参数
type SimpleSearchQuery struct {
	Index  []string
	IDs    []string
	From   *int
	Size   *int
	Pretty *bool
	Sort   []Sort
}

func (s *SimpleSearchQuery) getSize() *int {
	size := 500
	if s.Size == nil {
		return &size
	}
	return s.Size
}

// Sort 排序方式
// Field 要排序的字段
// Ascending 排序方式,true:正序 false:倒序
type Sort struct {
	Field     string
	Ascending bool
}

// QueryEntity 查询参数
// Key queryKey指定查询方式 Term Range ...
// Values QueryBody对象
type QueryEntity struct {
	Key    queryKey
	Values []QueryBody
}

// QueryBody 查询参数body值集合
// Key 作为查询条件的字段
// Value Term查询用的值
// Gte Range查询 >= 的值
// Lte Range查询 <= 的值
// Gt Range查询 > 的值
// Lt Range查询 < 的值
type QueryBody struct {
	Key   string
	Value any
	Gte   any
	Lte   any
	Gt    any
	Lt    any
}
