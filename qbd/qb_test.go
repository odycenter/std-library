package qbd_test

import (
	"fmt"
	"std-library/qbd"
	"testing"
)

func TestNewQueryBuilder(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	qb, _ := qbd.NewQueryBuilder("mysql")
	qb.Select("a.*", "b.a").
		From("table_a AS a").
		InnerJoin("table_b AS b").On("a.a = b.a").
		Where("a.a = ?").
		And("b.b >= ?").
		And("a.c").
		In(len(in)).
		GroupBy("b.a").
		Having("SUM(b.d) > 1000").
		Limit(1).
		Offset(0)
	fmt.Println(qb.String())
}
