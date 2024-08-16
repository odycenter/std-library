package module

import (
	"embed"
	"log"
	"log/slog"
	"os"
	app "std-library/app/conf"
	applog "std-library/app/log"
	"std-library/logs"
	"strings"
)

type SystemModule struct {
	Common
	EnvProperties map[string]embed.FS
}

func (m *SystemModule) Initialize() {
	if m.EnvProperties == nil || len(m.EnvProperties) == 0 {
		log.Panic("EnvProperties is empty")
	}

	handler := applog.NewHandler(os.Stdout)
	handler.SetLevel(slog.LevelDebug)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	m.LoadProperties(m.EnvProperties, "sys.properties")
	appName := m.RequiredProperty("core.app.name")
	logs.AppName = appName
	app.Name = appName
	if app.Local() {
		m.ModuleContext.PropertyManager.EnableLocalPropertyOverride(appName)
	}

	m.configureLog()
	m.configureCache()
	m.configureRedis()
	m.configureDB()
	m.configureMongo()
	m.configureKafka()
	m.configurePyroScope()
	m.configureMetric()
	m.configureGRPC()
	m.configureHTTP()
}

func (m *SystemModule) configureLog() {
	m.Log().DefaultLevel(m.Property("sys.log.level"))
	if app.Local() {
		slog.Info("Setting log level to DEBUG for local environment")
		m.Log().DefaultLevel(slog.LevelDebug.String())
	}
}

func (m *SystemModule) configureHTTP() {
	httpListen := m.Property("sys.http.listen")
	if httpListen == "" {
		return
	}

	m.Http().Listen(httpListen)

	allowCIDR := m.Property("sys.api.allowCIDR")
	if allowCIDR != "" {
		m.Http().AllowAPI(NewIPv4RangePropertyValueParser(allowCIDR).Parse())
	}

	m.Http().APIContent(&m.EnvProperties)
}

func (m *SystemModule) configureDB() {
	// TODO
}

func (m *SystemModule) configureMongo() {
	mongoUri := m.Property("sys.mongo.uri")
	if mongoUri != "" {
		m.Mongo().Uri(mongoUri)
	}
	user := m.Property("sys.mongo.user")
	if user != "" {
		m.Mongo().User(user)

	}
	password := m.Property("sys.mongo.password")
	if password != "" {
		m.Mongo().Password(password)
	}
	auth := m.Property("sys.mongo.auth")
	if strings.ToLower(auth) == "iam" {
		m.Mongo().IAMAuth()
	}
}

func (m *SystemModule) configurePyroScope() {
	pyroscopeUri := m.Property("sys.pyroscope.uri")
	if pyroscopeUri != "" {
		m.Pyroscope().Uri(pyroscopeUri)
	}
}

func (m *SystemModule) configureKafka() {
	kafkaUri := m.Property("sys.kafka.uri")
	if kafkaUri != "" {
		m.Kafka().Uri(kafkaUri)
	}
}

func (m *SystemModule) configureRedis() {
	redisHost := m.Property("sys.redis.host")
	if redisHost != "" {
		m.Redis().Host(redisHost)
	}
}

func (m *SystemModule) configureGRPC() {
	grpcListen := m.Property("sys.grpc.listen")
	if grpcListen != "" {
		m.Grpc().Listen(grpcListen)
	}
}

func (m *SystemModule) configureCache() {
	host := m.Property("sys.cache.host")
	if host != "" {
		if host == "local" {
			m.Cache().Local()
		} else {
			m.Cache().Redis(host)
		}
	}
}

func (m *SystemModule) configureMetric() {
	config := m.Metric()
	metricListen := m.Property("sys.metric.listen")
	if metricListen != "" {
		config.Listen(metricListen)
	}
}
