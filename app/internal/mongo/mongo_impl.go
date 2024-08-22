package internal_mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	app "std-library/app/conf"
	actionlog "std-library/app/log"
	"strings"
	"time"
)

type MongoImpl struct {
	name                          string
	uri                           string
	option                        *options.ClientOptions
	credential                    *options.Credential
	client                        *mongo.Client
	slowOperationThresholdInNanos int64
	initialized                   bool
}

func New(name string) *MongoImpl {
	option := options.Client()
	option.SetAppName(app.Name)
	option.SetMinPoolSize(10)
	option.SetMaxPoolSize(50)
	option.SetMaxConnecting(50)
	option.SetMaxConnIdleTime(30 * time.Minute)
	option.SetRetryReads(true)
	option.SetRetryWrites(true)
	option.SetConnectTimeout(5 * time.Second)
	impl := &MongoImpl{
		name:                          name,
		option:                        option,
		slowOperationThresholdInNanos: 1 * time.Second.Nanoseconds(),
	}
	impl.Timeout(120 * time.Second)
	return impl
}

func (m *MongoImpl) Execute(ctx context.Context) {
	if m.initialized {
		return
	}

	slog.InfoContext(ctx, fmt.Sprintf("mongo Initialize, name=%s", m.name))
	m.Initialize()
}

func (m *MongoImpl) Initialized() bool {
	return m.initialized
}

func (m *MongoImpl) Client() *mongo.Client {
	if !m.initialized {
		log.Fatalf("mongo is not initialized, name=%s", m.name)
	}
	return m.client
}

func (m *MongoImpl) Initialize() {
	if m.credential != nil {
		m.option.SetAuth(*m.credential)
	}
	m.option.SetMonitor(m.Monitor())
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, m.option)
	if err != nil {
		log.Panicf("[mongo] Connect to <%s> Failed, name: %s, err: %v\n", m.uri, m.name, err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("[mongo] Ping to <%s> Failed, name: %s, err: %v\n", m.uri, m.name, err)
	}
	m.client = client
	m.initialized = true
}

func (m *MongoImpl) Close() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		slog.Error("close mongo client error", "name", m.name, "uri", m.uri, "error", err)
	}
}

func (m *MongoImpl) Get() *mongo.Client {
	return m.client
}

func (m *MongoImpl) Uri(uri string) {
	m.uri = uri
	m.option.ApplyURI(uri)
}

func (m *MongoImpl) User(user string) {
	if m.credential == nil {
		m.credential = &options.Credential{}
	}
	m.credential.Username = user
}

func (m *MongoImpl) Password(password string) {
	if m.credential == nil {
		m.credential = &options.Credential{}
	}
	m.credential.Password = password
}

func (m *MongoImpl) IAMAuth() {
	if m.credential == nil {
		m.credential = &options.Credential{}
	}
	m.credential.AuthMechanism = "MONGODB-AWS"
	m.credential.AuthSource = "$external"
}

func (m *MongoImpl) Timeout(timeout time.Duration) {
	m.option.SetTimeout(timeout)
	m.option.SetSocketTimeout(timeout)
	m.option.SetServerSelectionTimeout(3 * timeout)
}

func (m *MongoImpl) PoolSize(minSize, maxSize uint64) {
	m.option.SetMinPoolSize(minSize)
	m.option.SetMaxPoolSize(maxSize)
}

func (m *MongoImpl) MaxConnecting(maxConnecting uint64) {
	m.option.SetMaxConnecting(maxConnecting)
}

func (m *MongoImpl) TLSConfig(conf *tls.Config) {
	if conf != nil {
		m.option.SetTLSConfig(conf)
	}
}

func (m *MongoImpl) SlowOperationThreshold(d time.Duration) {
	m.slowOperationThresholdInNanos = d.Nanoseconds()
}

func (m *MongoImpl) Monitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			elapsed := event.CommandFinishedEvent.Duration.Nanoseconds()
			actionlog.Stat(&ctx, "mongo."+strings.ToLower(event.CommandName), float64(elapsed))

			if elapsed > m.slowOperationThresholdInNanos {
				actionlog.Context(&ctx, "slow_operation", true)
				slog.WarnContext(ctx, fmt.Sprintf("[SLOW_OPERATION] slow %s, duration %v, db: %s", event.CommandName, event.CommandFinishedEvent.Duration, event.DatabaseName))
			}
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			elapsed := event.CommandFinishedEvent.Duration.Nanoseconds()
			actionlog.Stat(&ctx, "mongo."+strings.ToLower(event.CommandName), float64(elapsed))

			if elapsed > m.slowOperationThresholdInNanos {
				actionlog.Context(&ctx, "slow_operation", true)
				slog.WarnContext(ctx, fmt.Sprintf("[SLOW_OPERATION] slow %s, duration %v, db: %s", event.CommandName, event.CommandFinishedEvent.Duration, event.DatabaseName))
			}
		},
	}

}
