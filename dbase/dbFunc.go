package dbase

import (
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

// ErrorFieldsIllegal 查询参数错误
var ErrorFieldsIllegal = errors.New("[dbase]Query fields must be 0 or a multiple of 2")
var ErrSQLInjected = errors.New("[dbase]SQL inject")

// 校验filter是否符合规定
func verifyFields(fields []any) error {
	if len(fields)%2 != 0 {
		return ErrorFieldsIllegal
	}
	return nil
}

// List 用于返回db查询的一组数据，以传入数组ptr方式获取查询返回值
// in条件field key自行添加 __in, val为数组
func (d *DB) List(m any, i any, orders *[]string, cols *[]string, fields ...any) (rows int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	if cols != nil {
		return qs.AllWithCtx(d.getCtx(), i, *cols...)
	}
	return qs.AllWithCtx(d.getCtx(), i)
}

// ListRaw 用于返回db查询的一组数据，以[]orm.Params形式返回
// in条件field key自行添加 __in, val为数组
func (d *DB) ListRaw(m any, orders *[]string, fields ...any) (rows int64, data []orm.Params, err error) {
	err = verifyFields(fields)
	if err != nil {
		return 0, nil, err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	rows, err = qs.ValuesWithCtx(d.getCtx(), &data)
	return
}

// One 用于返回db查询List的第一条数据，以传入数组ptr方式获取查询返回值
func (d *DB) One(m any, i any, orders *[]string, fields ...any) (err error) {
	err = verifyFields(fields)
	if err != nil {
		return err
	}
	qs := d.QueryTable(m)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	if orders != nil {
		qs = qs.OrderBy(*orders...)
	}
	qs = qs.Limit(-1)
	err = qs.OneWithCtx(d.getCtx(), i)
	return
}

// Get 用于查询一条数据，以传入数组ptr方式获取查询返回值
func (d *DB) Get(m any, cols ...string) (err error) {
	return d.ReadWithCtx(d.getCtx(), m, cols...)
}

// InsertMulti 一次插入多条条数据
// perIns 单次插入数量
func (d *DB) InsertMulti(i any, perIns int) (id int64, err error) {
	return d.InsertMultiWithCtx(d.getCtx(), perIns, i)
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
func (d *DB) UpgradeFilter(i any, filters *map[string]any, values *orm.Params) (rows int64, err error) {
	qs := d.QueryTable(i)
	if filters != nil {
		for k, v := range *filters {
			qs = qs.Filter(k, v)
		}
	}
	return qs.UpdateWithCtx(d.getCtx(), *values)
}

// InsertOrUpdate 插入或更新一条数据
func (d *DB) InsertOrUpdate(i any, fields ...string) (rows int64, err error) {
	return d.InsertOrUpdateWithCtx(d.getCtx(), i, fields...)
}

// Delete 删除数据
func (d *DB) Delete(i any, fields ...string) (rows int64, err error) {
	return d.DeleteWithCtx(d.getCtx(), i, fields...)
}

// Count 数量
func (d *DB) Count(i any, fields ...any) (count int64, err error) {
	err = verifyFields(fields)
	if err != nil {
		return
	}
	qs := d.QueryTable(i)
	for i := 0; i < len(fields)/2; i++ {
		qs = qs.Filter(fields[i*2+0].(string), fields[i*2+1])
	}
	count, err = qs.CountWithCtx(d.getCtx())
	return
}

// Begin 创建事务
func (d *DB) Begin() (*TxOrm, error) {
	tx, err := d.Ormer.BeginWithCtx(d.getCtx())
	if err != nil {
		return nil, err
	}
	return &TxOrm{TxOrmer: tx}, nil
}

// Tx 创建回调事务
func (d *DB) Tx(fn func(ctx context.Context, tx *TxOrm) error) error {
	return d.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		return fn(ctx, &TxOrm{TxOrmer: txOrm})
	})
}
