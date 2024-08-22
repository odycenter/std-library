package internal_db

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/go-sql-driver/mysql"
)

const (
	defaultPoolMinSize       = 15
	defaultPoolMaxSize       = 40
	defaultConnMaxIdleTime   = 30 * time.Minute
	longTransactionThreshold = 5 * time.Second
)

var (
	once sync.Once
)

type DBImpl struct {
	name            string
	user            string
	password        string
	url             string
	alias           string
	region          string
	poolMinSize     int
	poolMaxSize     int
	connMaxIdleTime time.Duration
	authProvider    *AWSAuthProvider
	db              *sql.DB
	initialized     bool
}

func New(name string) *DBImpl {
	ConfigureLog()

	return &DBImpl{
		name:            name,
		poolMinSize:     defaultPoolMinSize,
		poolMaxSize:     defaultPoolMaxSize,
		connMaxIdleTime: defaultConnMaxIdleTime,
	}
}

func (d *DBImpl) Execute(ctx context.Context) {
	if d.initialized {
		return
	}

	slog.InfoContext(ctx, "Initializing database",
		"name", d.name,
		"alias", d.alias,
		"user", d.user,
		"url", d.url,
		"poolMinSize", d.poolMinSize,
		"poolMaxSize", d.poolMaxSize,
		"connMaxIdleTime", d.connMaxIdleTime,
	)

	d.initialize()
}

func (d *DBImpl) Initialized() bool {
	return d.initialized
}

func (d *DBImpl) initialize() {
	if d.initialized {
		return
	}

	config := d.authConfig()
	connector, err := mysql.NewConnector(config)
	if err != nil {
		log.Fatalf("DB connector create failed, name=%s, err=%v", d.name, err)
	}

	d.db = sql.OpenDB(connector)

	name := d.name
	if d.alias != "" {
		slog.Debug("Registering DB by alias", "name", d.name, "alias", d.alias)
		name = d.alias
	}

	err = orm.AddAliasWthDB(name, "mysql", d.db,
		orm.MaxIdleConnections(d.poolMinSize),
		orm.MaxOpenConnections(d.poolMaxSize),
		orm.ConnMaxIdletime(d.connMaxIdleTime),
		orm.ConnMaxLifetime(d.connMaxIdleTime*3))
	if err != nil {
		d.Close()
		log.Fatalf("register DB `%s`, err=%v", d.name, err)
	}

	d.initialized = true
}

func (d *DBImpl) authConfig() (config *mysql.Config) {
	config, err := mysql.ParseDSN(d.url)
	if err != nil {
		log.Fatalf("DB uri parse failed, name=%s uri=%s, err=%v", d.name, d.url, err)
	}
	config.User = d.user
	config.Params = map[string]string{
		"charset": "utf8mb4",
	}
	config.CheckConnLiveness = true
	slog.Info("before setting password to config", "dsn", config.FormatDSN())
	config.Passwd = d.password

	if iamUser(config.User) {
		if d.region == "" {
			log.Fatalf("IAM auth requires a non-empty region")
		}
		RegisterCerts()

		d.authProvider = &AWSAuthProvider{
			User:       config.User,
			DBEndpoint: config.Addr,
			Region:     d.region,
		}
		config.TLSConfig = "rds"
		config.AllowCleartextPasswords = true

		beforeConnect := mysql.BeforeConnect(func(ctx context.Context, cfg *mysql.Config) error {
			cfg.Passwd = d.authProvider.AccessToken()
			return nil
		})
		if err := config.Apply(beforeConnect); err != nil {
			log.Fatalf("apply config: %v", err)
		}
		slog.Info("Connecting to MySQL using IAM", "dsn", config.FormatDSN())
	}

	return config
}

func (d *DBImpl) Close() {
	if d.db == nil {
		return
	}
	err := d.db.Close()
	if err != nil {
		slog.Error("close db error", "name", d.name, "error", err)
	}
}

func (d *DBImpl) Url(url string) {
	d.url = url
}

func (d *DBImpl) User(user string)                { d.user = user }
func (d *DBImpl) Password(password string)        { d.password = password }
func (d *DBImpl) PoolSize(minSize, maxSize int)   { d.poolMinSize, d.poolMaxSize = minSize, maxSize }
func (d *DBImpl) Region(region string)            { d.region = region }
func (d *DBImpl) Alias(alias string)              { d.alias = alias }
func (d *DBImpl) ConnMaxIdleTime(t time.Duration) { d.connMaxIdleTime = t }

func ConfigureLog() {
	once.Do(func() {
		orm.Debug = true
		orm.DebugLog = orm.NewLog(&DoNothingWriter{})
		orm.LogFunc = logQuery
	})
}

func logQuery(query map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(query)+1)
	for k, v := range query {
		attrs = append(attrs, slog.Any(k, v))
	}

	level := slog.LevelDebug
	if flag, ok := query["flag"].(string); ok && flag == "FAIL" {
		level = slog.LevelError
	}
	if costTime, ok := query["cost_time"].(float64); ok && costTime >= float64(longTransactionThreshold.Milliseconds()) {
		if level != slog.LevelError {
			level = slog.LevelWarn
		}
		attrs = append(attrs, slog.Bool("slow_process", true))
	}
	slog.LogAttrs(context.Background(), level, "", attrs...)
	// TODO _sys/db api to change output level?
}

func iamUser(user string) bool {
	return strings.Contains(user, "_iam") || strings.Contains(user, "iam_")
}
