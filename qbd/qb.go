// Package qbd SQl Query Builder SQL语句构建器
package qbd

import "errors"

// QueryBuilder is the Query builder interface
type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	ForUpdate() QueryBuilder
	From(tables ...string) QueryBuilder
	InnerJoin(table string) QueryBuilder
	LeftJoin(table string) QueryBuilder
	RightJoin(table string) QueryBuilder
	On(cond string) QueryBuilder
	Where(cond string) QueryBuilder
	And(cond string) QueryBuilder
	Or(cond string) QueryBuilder
	In(vLen int) QueryBuilder
	InSQL(sql string) QueryBuilder
	InVals(vals ...string) QueryBuilder //Deprecated: 非SQL注入安全方法，请勿使用
	OrderBy(fields ...string) QueryBuilder
	Asc() QueryBuilder
	Desc() QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Having(cond string) QueryBuilder
	Update(tables ...string) QueryBuilder
	Set(kv ...string) QueryBuilder
	Delete(tables ...string) QueryBuilder
	InsertInto(table string, fields ...string) QueryBuilder
	Values(vals ...string) QueryBuilder
	Subquery(sub string, alias string) string
	String(multiplex ...bool) string
}

// NewQueryBuilder return the QueryBuilder
func NewQueryBuilder(driver string) (qb QueryBuilder, err error) {
	if driver == "mysql" {
		qb = new(MySQLQueryBuilder)
	} else if driver == "postgres" {
		qb = new(PostgresQueryBuilder)
	} else if driver == "sqlite" {
		err = errors.New("sqlite query builder is not supported yet")
	} else {
		err = errors.New("unknown driver for query builder")
	}
	return
}
