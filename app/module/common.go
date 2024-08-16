package module

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"std-library/app/async"
	app "std-library/app/conf"
	internal "std-library/app/internal/module"
)

type Common struct {
	ModuleContext *Context
	Env           string
}

func (c *Common) Load(module Module) {
	module.SetContext(c.ModuleContext)
	module.Initialize()
}

func (c *Common) SetContext(moduleContext *Context) {
	c.ModuleContext = moduleContext
}

func (c *Common) OnStartup(task async.Task) {
	c.ModuleContext.StartupHook.Add(task)
}

func (c *Common) OnShutdown(task async.Task) {
	c.ModuleContext.ShutdownHook.Add(internal.STAGE_5, func(ctx context.Context, timeoutInMs int64) {
		task.Execute(ctx)
	})
}

func (c *Common) LoadProperties(envFS map[string]embed.FS, propertiesFileName string) {
	c.Env = app.Env()
	slog.Info(fmt.Sprintf("loadProperties by env: %s, propertiesFileName: %s", c.Env, propertiesFileName))
	files, ok := envFS[c.Env]
	if !ok {
		log.Fatal("loadProperties error! Invalid environment: " + c.Env)
	}

	defaultFS, ok := envFS[""]
	if !ok {
		log.Fatal("loadProperties error! Default environment not found!")
	}
	c.LoadPropertiesByFS(files, propertiesFileName, defaultFS)
}

func (c *Common) LoadPropertiesByFS(properties embed.FS, propertyFile string, defaultFS embed.FS) {
	f, err := properties.Open(propertyFile)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(fmt.Sprintf("propertyFile not found!  load default properties, env: %v, fileName: %s ", app.Env(), propertyFile))
			f, err = defaultFS.Open(propertyFile)
		}

		if err != nil {
			log.Fatal("propertyFile not found!", propertyFile, err)
		}
	}
	defer f.Close()
	c.ModuleContext.PropertyManager.LoadProperties(f)
}

func (c *Common) Property(key string) string {
	return c.ModuleContext.Property(key)
}

func (c *Common) RequiredProperty(key string) string {
	v := c.ModuleContext.Property(key)
	if v == "" {
		panic("Required property not found: " + key)
	}
	return v
}

func (c *Common) Kafka(name ...string) *KafkaConfig {
	return c.ModuleContext.Config(configName("kafka", name...), func() Config { return &KafkaConfig{} }).(*KafkaConfig)
}

func (c *Common) Schedule() *SchedulerConfig {
	return c.ModuleContext.Config("scheduler", func() Config { return &SchedulerConfig{} }).(*SchedulerConfig)
}

func (c *Common) Http() *HTTPConfig {
	return c.ModuleContext.Config("http", func() Config { return &HTTPConfig{} }).(*HTTPConfig)
}

func (c *Common) Grpc() *GrpcServerConfig {
	return c.ModuleContext.Config("grpc", func() Config { return &GrpcServerConfig{} }).(*GrpcServerConfig)
}

func (c *Common) Cache(name ...string) *CacheConfig {
	return c.ModuleContext.Config(configName("cache", name...), func() Config { return &CacheConfig{} }).(*CacheConfig)
}

func (c *Common) Redis(name ...string) *RedisConfig {
	return c.ModuleContext.Config(configName("redis", name...), func() Config { return &RedisConfig{} }).(*RedisConfig)
}

func (c *Common) Pyroscope() *PyroscopeConfig {
	return c.ModuleContext.Config("pyroscope", func() Config { return &PyroscopeConfig{} }).(*PyroscopeConfig)
}

func (c *Common) Mongo(name ...string) *MongoConfig {
	return c.ModuleContext.Config(configName("mongo", name...), func() Config { return &MongoConfig{} }).(*MongoConfig)
}

func (c *Common) Metric() *MetricConfig {
	return c.ModuleContext.Config("metric", func() Config { return &MetricConfig{} }).(*MetricConfig)
}

func (c *Common) Log() *LogConfig {
	return c.ModuleContext.Config("log", func() Config { return &LogConfig{} }).(*LogConfig)
}

func configName(prefix string, name ...string) string {
	cname := prefix
	if len(name) > 0 && name[0] != "" {
		cname = cname + ":" + name[0]
	}
	return cname
}
