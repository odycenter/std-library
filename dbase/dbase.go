// Package dbase 基于beego的orm实现封装的mySQL操作
package dbase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"std-library/dbase/cloud"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

// Opt 配置信息
type Opt struct {
	AliasName      string
	ServiceName    []string            //default []string,if have value(s),this option will only be used for this(these) service(s)--implement by init
	DriverName     string              // driver name
	DriverTyp      orm.DriverType      // driver type
	Region         string              // region
	Host           string              //host <must>
	Port           string              //default 3306
	User           string              //username <must>
	Password       string              //password
	DBName         string              //connect DB name <must>
	SslMode        string              //default disable
	TimeZone       *time.Location      //default local
	MaxIdleConnes  int                 //default 10
	MaxOpenConnes  int                 //default 30
	IsolationLevel string              //default ""
	OrmDebug       bool                //is open orm debug mode, output the SQl Exec
	OrmLogFunc     func(query *OrmLog) //when open orm debug mode,log func can be specified
}

type OrmLog struct {
	HostName string  `json:"host_name"`
	CostTime float64 `json:"cost_time"`
	Flag     string  `json:"flag"`
	Sql      string  `json:"sql"`
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

func (opt *Opt) getLogFunc() func(map[string]interface{}) {
	if opt.OrmLogFunc != nil {
		return func(m map[string]interface{}) {
			costTime, ok := m["cost_time"].(float64)
			if !ok {
			}
			flag, ok := m["flag"].(string)
			if !ok {
			}
			sql, ok := m["sql"].(string)
			if !ok {
			}
			hostName, _ := os.Hostname()
			opt.OrmLogFunc(&OrmLog{
				HostName: hostName,
				CostTime: costTime,
				Flag:     flag,
				Sql:      sql,
			})
		}
	}
	return nil
}

// Init 初始化orm数据库连接池
func Init(opts ...*Opt) {
	for _, opt := range opts {
		err := orm.RegisterDriver(opt.getDriverName(), opt.getDriverTyp())
		if err != nil {
			log.Panic(err)
		}
		forceRDSIAM := os.Getenv("FORCE_RDS_IAM_USER")
		if forceRDSIAM == "true" || opt.User == "cloud_iam" {
			RegisterDataBase(opt)
		} else {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", opt.User, opt.Password, opt.Host, opt.getPort(), opt.DBName)
			err = orm.RegisterDataBase(opt.getAliasName(), opt.getDriverName(), dsn,
				orm.MaxIdleConnections(opt.getMaxIdleConnes()),
				orm.MaxOpenConnections(opt.getMaxOpenConnes()),
				orm.ConnMaxIdletime(time.Hour*1),
				orm.ConnMaxLifetime(time.Hour*2))
			if err != nil {
				log.Panic(err)
			}
		}
		err = orm.SetDataBaseTZ(opt.getAliasName(), opt.getTimeZone())
		if err != nil {
			log.Panic(err)
		}
		orm.BootStrap()
		if opt.OrmDebug {
			orm.Debug = opt.OrmDebug
			if opt.OrmLogFunc != nil {
				orm.LogFunc = opt.getLogFunc()
			}
		}
	}
}

func RegisterDataBase(opt *Opt) {
	authProvider := cloud.AWSAuthProvider{
		DBEndpoint: opt.Host + ":" + opt.getPort(),
		Region:     opt.Region,
		DBName:     opt.DBName,
	}
	authProvider.Register()
	dsn := authProvider.DataSourceName()
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}
	cfg.CheckConnLiveness = true
	if opt.IsolationLevel != "" {
		if cfg.Params == nil {
			cfg.Params = map[string]string{"tx_isolation": opt.IsolationLevel}
		} else {
			cfg.Params["tx_isolation"] = opt.IsolationLevel
		}
	}
	beforeConnect := mysql.BeforeConnect(func(ctx context.Context, cfg *mysql.Config) error {
		if cfg.User == "cloud_iam" {
			cfg.Passwd = authProvider.AccessToken()
		}
		return nil
	})
	err1 := cfg.Apply(beforeConnect)
	if err1 != nil {
		panic(err1)
	}
	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		panic(err)
	}

	db := sql.OpenDB(connector)
	if err != nil {
		if db != nil {
			db.Close()
		}
		err = fmt.Errorf("register db `%s`, %s", opt.getAliasName(), err.Error())
		panic(err)
	}
	err = orm.AddAliasWthDB(opt.getAliasName(), opt.getDriverName(), db,
		orm.MaxIdleConnections(opt.getMaxIdleConnes()),
		orm.MaxOpenConnections(opt.getMaxOpenConnes()),
		orm.ConnMaxIdletime(time.Hour*1),
		orm.ConnMaxLifetime(time.Hour*2))
	if err != nil {
		log.Panic(err)
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
	if d.ctx == nil {
		return context.Background()
	}
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

func GetDB(aliasName ...string) (*sql.DB, error) {
	name := "default"
	if len(aliasName) != 0 && aliasName[0] != "" {
		name = aliasName[0]
	}
	return orm.GetDB(name)
}
