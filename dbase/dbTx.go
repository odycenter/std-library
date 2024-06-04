package dbase

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
)

type TxOrm struct {
	orm.TxOrmer
	ctx context.Context
}

// WithCtx 传入自定义context
func (tx *TxOrm) WithCtx(ctx context.Context) *TxOrm {
	tx.ctx = ctx
	return tx
}

// WithCtx 传入自定义context
func (tx *TxOrm) getCtx() context.Context {
	if tx.ctx == nil {
		return context.Background()
	}
	return tx.ctx
}

// UpgradeFilter 按字段更新数据
// filter 条件
/* values Set = v 条件
orm.Params{
    "name": "astaxie",
}
或
orm.Params{
    "nums": orm.ColValue(orm.ColAdd, 100),
}
ColAdd      // 加
ColMinus    // 减
ColMultiply // 乘
ColExcept   // 除
*/
func (tx *TxOrm) UpgradeFilter(i interface{}, filters *map[string]interface{}, values *orm.Params) (rows int64, err error) {
	qs := tx.QueryTable(i)
	if filters != nil {
		for k, v := range *filters {
			qs = qs.Filter(k, v)
		}
	}
	return qs.Update(*values)
}

// List 用于返回db查询的一组数据，以传入数组ptr方式获取查询返回值
// in条件field key自行添加 __in, val为数组
func (tx *TxOrm) List(m any, i any, orders *[]string, cols *[]string, fields ...any) (rows int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, err
	}
	qs := tx.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	if cols != nil {
		return qs.AllWithCtx(tx.getCtx(), i, *cols...)
	}
	return qs.AllWithCtx(tx.getCtx(), i)
}

// ListRaw 用于返回db查询的一组数据，以[]orm.Params形式返回
// in条件field key自行添加 __in, val为数组
func (tx *TxOrm) ListRaw(m any, orders *[]string, fields ...any) (rows int64, data []orm.Params, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, nil, err
	}
	qs := tx.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	rows, err = qs.ValuesWithCtx(tx.getCtx(), &data)
	return
}

// One 用于返回db查询List的第一条数据，以传入数组ptr方式获取查询返回值
func (tx *TxOrm) One(m any, i any, orders *[]string, fields ...any) (err error) {
	err = verifyFields(fields)
	if err != nil {
		return err
	}
	qs := tx.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	err = qs.OneWithCtx(tx.getCtx(), i)
	return
}

// Get 用于查询一条数据，以传入数组ptr方式获取查询返回值
func (tx *TxOrm) Get(m any, cols ...string) (err error) {
	return tx.ReadWithCtx(tx.getCtx(), m, cols...)
}

// InsertOrUpdate 插入或更新一条数据
func (tx *TxOrm) InsertOrUpdate(i any, fields ...string) (rows int64, err error) {
	return tx.InsertOrUpdateWithCtx(tx.getCtx(), i, fields...)
}

// Delete 删除数据
func (tx *TxOrm) Delete(i any, fields ...string) (rows int64, err error) {
	return tx.DeleteWithCtx(tx.getCtx(), i, fields...)
}

// DeleteMany 按条件删除多条数据
func (tx *TxOrm) DeleteMany(i any, filters ...any) (rows int64, err error) {
	q := tx.QueryTableWithCtx(tx.getCtx(), i)
	for i := 0; i < len(filters)/2; i++ {
		q = q.Filter(filters[i*2].(string), filters[i*2+1])
	}
	return q.Delete()
}

// Count 数量
func (tx *TxOrm) Count(i any, fields ...any) (count int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return
	}
	qs := tx.QueryTable(i)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	count, err = qs.CountWithCtx(tx.getCtx())
	return
}
