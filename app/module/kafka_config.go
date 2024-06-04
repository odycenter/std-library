package module

import (
	"context"
	"log"
	"runtime"
	app "std-library/app/conf"
	internal "std-library/app/internal/module"
	"std-library/app/kafka"
	"std-library/logs"
	"sync"
)

type KafkaConfig struct {
	name          string
	groupId       string
	moduleContext *Context
	uriString     string
	uri           []string
	poolSize      int
	m             map[string]*kafka.MessageListener
	mu            sync.RWMutex
	handlerAdded  bool
}

func (c *KafkaConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext
	c.groupId = app.Name
	c.poolSize = runtime.NumCPU() * 4
	if c.poolSize > 4 {
		c.poolSize = 4
	}
	logs.Info("kafka consumer default poolSize: ", c.poolSize)
	c.m = make(map[string]*kafka.MessageListener)
	c.moduleContext.StartupHook.Add(c)
	c.moduleContext.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		c.stop(ctx, timeoutInMs)
	})
}

func (c *KafkaConfig) Validate() {
	if !c.handlerAdded {
		log.Fatalf("kafka is configured, but no handler added, please remove unnecessary config, name=" + c.name)
	}
	if len(c.uri) == 0 {
		log.Fatalf("kafka uri is not configured, name=" + c.name)
	}
}

func (c *KafkaConfig) Uri(uri string) {
	if c.uriString != "" {
		log.Fatalf("kafka uri is already configured, uri=%s, previous=%s", uri, c.uri)
	}
	u := kafka.Uri{}
	u.Uri(uri)
	c.uri = u.Parse()
	c.moduleContext.Probe.AddHostURI(c.uri[0])
}

// GroupId by default use AppName as consumer group
// use "${service-name}-${label}" to allow same service to be deployed for multitenancy
func (c *KafkaConfig) GroupId(groupId string) {
	c.groupId = groupId
}

func (c *KafkaConfig) DefaultPoolSize(size int) {
	c.poolSize = size
}

// Subscribe that use app.Name as groupID.
func (c *KafkaConfig) Subscribe(topic string, handler kafka.MessageHandler, poolSize ...int) {
	c.subscribe(topic, kafka.FirstOffset, handler, poolSize...)
}

func (c *KafkaConfig) SubscribeByOffset(topic string, startOffset int64, handler kafka.MessageHandler, poolSize ...int) {
	c.subscribe(topic, startOffset, handler, poolSize...)
}

func (c *KafkaConfig) subscribe(topic string, startOffset int64, handler kafka.MessageHandler, poolSize ...int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.m[topic]
	if ok {
		log.Fatalf("topic is already subscribed, topic: %s ", topic)
	}

	listener := kafka.MessageListener{
		Handler: handler,
	}
	if len(poolSize) > 0 && poolSize[0] > 0 {
		listener.SetPoolSize(poolSize[0])
	} else {
		listener.SetPoolSize(c.poolSize)
	}
	opt := &kafka.SubscribeOption{
		Topic:       topic,
		StartOffset: startOffset,
	}
	listener.Initialize(opt)

	c.m[opt.Topic] = &listener
	c.handlerAdded = true
}

func (c *KafkaConfig) start(ctx context.Context) {
	c.mu.RLock()
	for _, listener := range c.m {
		listener.Opt.GroupId = c.groupId
		listener.Opt.Brokers = c.uri
		listener.Start(ctx)
	}
	c.mu.RUnlock()
}

func (c *KafkaConfig) stop(ctx context.Context, timeoutInMs int64) {
	var wg sync.WaitGroup
	c.mu.RLock()
	for _, listener := range c.m {
		wg.Add(1)
		go func(listener *kafka.MessageListener) {
			defer wg.Done()
			listener.AwaitTermination(ctx, timeoutInMs)
		}(listener)
	}
	c.mu.RUnlock()
	wg.Wait()
}

func (c *KafkaConfig) Execute(ctx context.Context) {
	c.start(ctx)
}
