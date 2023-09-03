package dbase

import (
	"github.com/beego/beego/v2/client/orm"
)

type TxOrm struct {
	orm.TxOrmer
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
