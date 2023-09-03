// Package dbase 基于beego的orm实现封装的mySQL操作
package dbase

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

// Opt 配置信息
type Opt struct {
	AliasName         string
	ServiceName       []string //default []string,if have value(s),this option will only be used for this(these) service(s)--implement by init
	DriverName        string
	DriverTyp         orm.DriverType
	Host              string //host <must>
	Port              string //default 3306
	User              string //username <must>
	Password          string
	DBName            string         //connect DB name <must>
	SslMode           string         //default disable
	TimeZone          *time.Location //default local
	MaxIdleConnes     int            //default 10
	MaxOpenConnes     int            //default 30
	MaxLifeTimeConnes time.Duration  //default 3600
	SyncDB            bool           //is need orm auto sync DB struct
}

func (opt *Opt) getDriverName() string {
	if opt.DriverName == "" {
		return "mysql"
	}
	return opt.DriverName
}

func (opt *Opt) getDriverTyp() orm.DriverType {
	if opt.DriverTyp == 0 {
		return orm.DRMySQL
	}
	return opt.DriverTyp
}

func (opt *Opt) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Opt) getPort() string {
	if opt.Port == "" {
		return "3306"
	}
	return opt.Port
}

func (opt *Opt) getMaxIdleConnes() int {
	if opt.MaxOpenConnes == 0 {
		return 10
	}
	return opt.MaxOpenConnes
}

func (opt *Opt) getMaxOpenConnes() int {
	if opt.MaxOpenConnes == 0 {
		return 30
	}
	return opt.MaxOpenConnes
}

func (opt *Opt) getMaxLifeTimeConnes() time.Duration {
	if opt.MaxLifeTimeConnes == 0 {
		return 3600 * time.Second
	}
	return opt.MaxLifeTimeConnes
}

func (opt *Opt) getTimeZone() *time.Location {
	if opt.TimeZone == nil {
		return time.Local
	}
	return opt.TimeZone
}

func (opt *Opt) getSslMode() string {
	if opt.SslMode == "" {
		return "disable"
	}
	return opt.SslMode
}

// Init 初始化orm数据库连接池
func Init(opts ...*Opt) {
	for _, opt := range opts {
		err := orm.RegisterDriver(opt.getDriverName(), opt.getDriverTyp())
		if err != nil {
			log.Panic(err)
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", opt.User, opt.Password, opt.Host, opt.getPort(), opt.DBName)
		aliasName := opt.getAliasName()
		err = orm.RegisterDataBase(aliasName, opt.getDriverName(), dsn,
			orm.MaxIdleConnections(opt.getMaxIdleConnes()),
			orm.MaxOpenConnections(opt.getMaxOpenConnes()),
			orm.ConnMaxLifetime(opt.getMaxLifeTimeConnes()))
		if err != nil {
			log.Panic(err)
		}
		err = orm.SetDataBaseTZ(aliasName, opt.getTimeZone())
		if err != nil {
			log.Panic(err)
		}
		if opt.SyncDB {
			err = orm.RunSyncdb(aliasName, false, true)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

// DB Database结构体
type DB struct {
	orm.Ormer
	ctx context.Context
}

// WithCtx 传入自定义context
func (d *DB) WithCtx(ctx context.Context) *DB {
	d.ctx = ctx
	return d
}

// WithCtx 传入自定义context
func (d *DB) getCtx() context.Context {
	return d.ctx
}

// Orm 获取一个orm，orm本身有连接管理
func Orm(aliasName ...string) *DB {
	name := "default"
	if len(aliasName) != 0 && aliasName[0] != "" {
		name = aliasName[0]
	}
	return &DB{orm.NewOrmUsingDB(name), context.Background()}
}

// IsDuplicate 是否为重复键报错
func IsDuplicate(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry")
}
