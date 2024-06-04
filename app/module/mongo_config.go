package module

import (
	"context"
	"crypto/tls"
	"log"
	internal "std-library/app/internal/module"
	internalmongo "std-library/app/internal/mongo"
	"time"
)

type MongoConfig struct {
	name          string
	moduleContext *Context
	uri           string
	mongoImpl     *internalmongo.MongoImpl
}

func (c *MongoConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext
	c.mongoImpl = internalmongo.New(name)
	c.moduleContext.StartupHook.Initialize = append(c.moduleContext.StartupHook.Initialize, c.mongoImpl)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_6, func(ctx context.Context, timeoutInMs int64) {
		c.mongoImpl.Close()
	})
}

func (c *MongoConfig) ForceEarlyStart() {
	c.Validate()
	c.mongoImpl.Execute(context.Background())
}

func (c *MongoConfig) Validate() {
	if c.uri == "" {
		log.Fatalf("mongo uri must be configured, name=" + c.name)
	}
}

func (c *MongoConfig) Uri(uri string) {
	if c.uri != "" {
		log.Fatalf("mongo uri is already configured, uri=%s, previous=%s", uri, c.uri)
	}
	c.uri = uri
	c.mongoImpl.Uri(uri)
	c.mongoImpl.TLSConfig(&tls.Config{InsecureSkipVerify: true})
}

func (c *MongoConfig) User(user string) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set user! name=" + c.name)
	}
	c.mongoImpl.User(user)
}

func (c *MongoConfig) Password(password string) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set password! name=" + c.name)
	}
	c.mongoImpl.Password(password)
}

func (c *MongoConfig) IAMAuth() {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set password! name=" + c.name)
	}
	c.mongoImpl.IAMAuth()
}

func (c *MongoConfig) Timeout(timeout time.Duration) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set timeout! name=" + c.name)
	}
	c.mongoImpl.Timeout(timeout)
}

func (c *MongoConfig) PoolSize(minSize, maxSize uint64) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set pool size! name=" + c.name)
	}
	c.mongoImpl.PoolSize(minSize, maxSize)
}

func (c *MongoConfig) MaxConnecting(maxConnecting uint64) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set maxConnecting! name=" + c.name)
	}
	c.mongoImpl.MaxConnecting(maxConnecting)
}

func (c *MongoConfig) TLSConfig(conf *tls.Config) {
	if c.mongoImpl.Initialized() {
		log.Fatalf("mongo is already initialized, can not set TLSConfig! name=" + c.name)
	}
	c.mongoImpl.TLSConfig(conf)
}
