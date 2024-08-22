package module

import (
	"context"
	"log"
	internalDB "std-library/app/internal/db"
	internal "std-library/app/internal/module"
	"strconv"
	"strings"
)

type DBConfig struct {
	name          string
	moduleContext *Context
	dbImpl        *internalDB.DBImpl
	url           string
}

func (c *DBConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext
	c.dbImpl = internalDB.New(name)
	c.moduleContext.StartupHook.Initialize = append(c.moduleContext.StartupHook.Initialize, c.dbImpl)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_6, func(ctx context.Context, timeoutInMs int64) {
		c.dbImpl.Close()
	})
}

func (c *DBConfig) Validate() {
	if c.url == "" {
		log.Fatalf("db url must be configured, name=" + c.name)
	}
}

func (c *DBConfig) ForceEarlyStart() {
	c.Validate()
	c.dbImpl.Execute(context.Background())
}

func (c *DBConfig) Url(url string) *DBConfig {
	if c.url != "" {
		log.Fatalf("DB uri is already configured, name=%s uri=%s, previous=%s", c.name, url, c.url)
	}
	c.url = url
	c.dbImpl.Url(url)
	return c
}

func (c *DBConfig) User(user string) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set user! name=%s", c.name)
	}

	c.dbImpl.User(user)
	return c
}

func (c *DBConfig) Password(password string) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set password! name=%s", c.name)
	}

	c.dbImpl.Password(password)
	return c
}

func (c *DBConfig) PoolSize(minSize int, maxSize int) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set pool size! name=%s", c.name)
	}

	c.dbImpl.PoolSize(minSize, maxSize)
	return c
}

func (c *DBConfig) PoolSizeString(setting string) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set pool size! name=%s", c.name)
	}

	parts := strings.Split(setting, ",")
	if len(parts) != 2 {
		log.Fatalf("Invalid pool size setting: %s. Expected format: minSize,maxSize", setting)
	}

	minSize, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		log.Fatalf("Invalid minSize in pool size setting: %s", parts[0])
	}

	maxSize, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		log.Fatalf("Invalid maxSize in pool size setting: %s", parts[1])
	}

	if minSize > maxSize {
		minSize, maxSize = maxSize, minSize
	}

	c.dbImpl.PoolSize(minSize, maxSize)
	return c
}

func (c *DBConfig) Region(region string) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set region! name=%s", c.name)
	}

	c.dbImpl.Region(region)
	return c
}

func (c *DBConfig) Alias(alias string) *DBConfig {
	if c.dbImpl.Initialized() {
		log.Fatalf("DB is already initialized, can not set alias! name=%s", c.name)
	}

	c.dbImpl.Alias(alias)
	return c
}
