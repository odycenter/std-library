package dbase_test

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
	"github.com/odycenter/std-library/dbase"
	"log"
	"testing"
)

type Element struct {
	A string `orm:"pk;index"`
	B int
}

func (*Element) TableName() string {
	return "element"
}

func TestDB(t *testing.T) {
	orm.RegisterModel(new(Element))
	dbase.Init(&dbase.Opt{
		DriverName:    "mysql",
		DriverTyp:     orm.DRMySQL,
		Host:          "127.0.0.1",
		Port:          "3307",
		DBName:        "test",
		User:          "root",
		Password:      "d@5xo4MkCYP&7qmB#6?4AAh$8trgo8ds",
		MaxIdleConnes: 1,
		MaxOpenConnes: 20,
	})
	o := dbase.Orm()
	var e = Element{"AAAAA", 123}
	_, err := o.Insert(&e)
	if err != nil && !dbase.IsDuplicate(err) {
		log.Println(err)
		return
	}
	e = Element{A: "AAAAA"}
	err = o.Get(&e, "A")
	if err != nil {
		log.Println(err)
		return
	}
	e.B = 222
	_, err = o.Update(&e, "B")
	if err != nil {
		log.Println(err)
		return
	}
	//_, err = o.Delete(&e, "A")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	/*----------------------------------------------------------------*/
	//orm事务使用示例1
	txOrm, _ := o.Begin()
	e.B = 1
	_, err = txOrm.Update(&e, "B")
	if err != nil {
		//回滚
		_ = txOrm.Rollback()
		log.Println(err)
		return
	}
	e.B = 999
	_, err = txOrm.Update(&e, "B")
	if err != nil {
		//回滚
		_ = txOrm.Rollback()
		log.Println(err)
		return
	}
	//提交
	err = txOrm.Commit()
	if err != nil {
		log.Println("commit failed ", err.Error())
		return
	}
	/*----------------------------------------------------------------*/
	//orm事务使用示例2
	txOrm, _ = o.Begin()
	defer func() {
		//整体回滚--
		err := txOrm.RollbackUnlessCommit()
		if err != nil {
			log.Println("rollback failed ", err.Error())
		}
	}()
	e.B = 1
	_, err = txOrm.Update(&e, "B")
	if err != nil {
		log.Println(err)
		return
	}
	e.B = 888
	_, err = txOrm.Update(&e, "B")
	if err != nil {
		log.Println(err)
		return
	}
	//提交
	err = txOrm.Commit()
	if err != nil {
		log.Println("commit failed ", err.Error())
		return
	}
	/*----------------------------------------------------------------*/
	//orm事务使用示例3
	err = o.Tx(func(ctx context.Context, txOrm *dbase.TxOrm) error {
		_, err = txOrm.Update(&e, "B")
		if err != nil {
			log.Println(err)
			return err
		}
		e.B = 777
		_, err = txOrm.Update(&e, "B")
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("commit failed ", err.Error())
		return
	}
}

func TestFilter(t *testing.T) {
	orm.RegisterModel(new(Element))
	dbase.Init(&dbase.Opt{
		DriverName:    "mysql",
		DriverTyp:     orm.DRMySQL,
		Host:          "127.0.0.1",
		Port:          "3307",
		DBName:        "test",
		User:          "root",
		Password:      "d@5xo4MkCYP&7qmB#6?4AAh$8trgo8ds",
		MaxIdleConnes: 1,
		MaxOpenConnes: 20,
	})
	o := dbase.Orm()
	var e = Element{}
	filter := o.Filter(&e, "A", "AAAAA")
	_, err := filter.UpdateWithCtx(context.TODO(), orm.Params{"B": 666})
	if err != nil {
		log.Println("Update failed ", err.Error())
		return
	}
}

func TestFilterRaw(t *testing.T) {
	orm.RegisterModel(new(Element))
	dbase.Init(&dbase.Opt{
		DriverName:    "mysql",
		DriverTyp:     orm.DRMySQL,
		Host:          "127.0.0.1",
		Port:          "3307",
		DBName:        "test",
		User:          "root",
		Password:      "d@5xo4MkCYP&7qmB#6?4AAh$8trgo8ds",
		MaxIdleConnes: 1,
		MaxOpenConnes: 20,
	})
	o := dbase.Orm()
	var e = Element{}
	filter := o.FilterRaw(&e, "A", "= 'AAAAA'")
	_, err := filter.UpdateWithCtx(context.TODO(), orm.Params{"B": 777})
	if err != nil {
		log.Println("Update failed ", err.Error())
		return
	}
}
