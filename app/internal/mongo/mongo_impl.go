package internal_mongo

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	app "std-library/app/conf"
	"std-library/logs"
	mongomigration "std-library/mongo"
	"time"
)

type MongoImpl struct {
	name        string
	uri         string
	option      *options.ClientOptions
	credential  *options.Credential
	client      *mongo.Client
	initialized bool
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
		name:   name,
		option: option,
	}
	impl.Timeout(120 * time.Second)
	return impl
}

func (m *MongoImpl) Execute(_ context.Context) {
	if m.initialized {
		return
	}

	logs.Debug("mongo Initialize, name=%s", m.name)
	m.Initialize()
	mongomigration.InitMigration(m.client)
}

func (m *MongoImpl) Initialized() bool {
	return m.initialized
}

func (m *MongoImpl) Initialize() {
	if m.credential != nil {
		m.option.SetAuth(*m.credential)
	}
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
		logs.Error("close mongo client error, name=%s, uri=%s, err=%v", m.name, m.uri, err)
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
