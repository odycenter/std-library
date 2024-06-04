package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	appkafka "std-library/app/kafka"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/logs"
	"sync"
	"time"
)

type consumersMap struct {
	sync.RWMutex
	m map[string][]*kafka.Reader
}

func init() {
	consumers.init()
}

func (cm *consumersMap) init() {
	cm.Lock()
	cm.m = make(map[string][]*kafka.Reader)
	cm.Unlock()
}

func (cm *consumersMap) load(k string) ([]*kafka.Reader, bool) {
	cm.RLock()
	defer cm.RUnlock()
	v, ok := cm.m[k]
	if ok {
		return v, true
	}
	return nil, false
}

func (cm *consumersMap) store(k string, v []*kafka.Reader) {
	cm.Lock()
	cm.m[k] = v
	cm.Unlock()
}

func (cm *consumersMap) del(k string) {
	cm.Lock()
	delete(cm.m, k)
	cm.Unlock()
}

type ConsumerOption struct {
	AliasName    string        `json:"AliasName"`    //别名
	Count        int           `json:"Count"`        //consumer启动数量
	BrokersAddrs []string      `json:"BrokersAddrs"` //Kafka单例、集群地址
	GroupID      string        `json:"GroupID"`      //消费者组ID
	Topics       []string      `json:"Topics"`       //订阅主题
	MaxBytes     int           `json:"MaxBytes"`     //消费者可接受的最大批量大小
	MinBytes     int           `json:"MinBytes"`     //消费者可接受的最小批量大小
	MaxWait      time.Duration `json:"MaxWait"`      //最大数据拉取等待时间
	ManualCommit bool          `json:"ManualCommit"` //是否手动提交 true手动提交 false自动提交
	Offset       int64         `json:"Offset"`       //起始偏移位
	OnReceive    OnReceive     `json:"-"`            //信息消费回调
	OnError      OnError       `json:"-"`            //消费过程中的错误
}

func (opt *ConsumerOption) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *ConsumerOption) getOffset() int64 {
	if opt.Offset == 0 {
		return LastOffset
	}
	return opt.Offset
}

func (opt *ConsumerOption) getMaxBytes() int {
	if opt.MaxBytes == 0 {
		return 10e5 // 10MB
	}
	return opt.MaxBytes
}

func (opt *ConsumerOption) getMaxWait() time.Duration {
	if opt.MaxWait == 0 {
		return 10 * time.Second
	}
	return opt.MaxWait
}

func (opt *ConsumerOption) getMinBytes() int {
	if opt.MinBytes == 0 {
		return 1
	}
	return opt.MinBytes
}

func (opt *ConsumerOption) getCount() int {
	if opt.Count == 0 {
		return 1 // 10MB
	}
	return opt.Count
}

// WithOnReceive 添加消息处理回调
func (opt *ConsumerOption) WithOnReceive(fn OnReceive) *ConsumerOption {
	opt.OnReceive = fn
	return opt
}

func (opt *ConsumerOption) getOnReceive() onReceive {
	return func(ctx context.Context, msg Message) {
		if opt.OnReceive == nil {
			return
		}
		opt.OnReceive(ctx, msg)
	}
}

// WithOnError 添加消息错误处理回调
func (opt *ConsumerOption) WithOnError(fn OnError) {
	opt.OnError = fn
}

func (opt *ConsumerOption) getOnError() OnError {
	if opt.OnError == nil {
		return func(err error) {}
	}
	return opt.OnError
}

// Deprecated: use Subscribe or SubscribeByOption instead, and do not use current Topic to subscribe.
// NewConsumer 使用ConsumerOption创建消费者
func NewConsumer(ctx context.Context, opt *ConsumerOption) {
	if _, ok := consumers.load(opt.getAliasName()); ok {
		log.Printf("kafka consumer <%s> already exists", opt.getAliasName())
		return
	}
	kafkaConfig := kafka.ReaderConfig{
		Brokers:     opt.BrokersAddrs,
		GroupID:     opt.GroupID,
		GroupTopics: opt.Topics,
		MaxBytes:    opt.getMaxBytes(), // 1MB
		MinBytes:    opt.getMinBytes(), //1
		MaxWait:     opt.getMaxWait(),  //def:10s
	}
	var readers []*kafka.Reader
	for i := 0; i < opt.Count; i++ {
		kafkaReader := kafka.NewReader(kafkaConfig)
		manualCommit := opt.ManualCommit
		go func(ctx context.Context, reader *kafka.Reader, opt *ConsumerOption) {
			do := opt.getOnReceive()
			var msg kafka.Message
			var err error
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if manualCommit {
						msg, err = kafkaReader.FetchMessage(ctx)
					} else {
						msg, err = kafkaReader.ReadMessage(ctx)
					}
					if err != nil {
						opt.getOnError()(err)
						continue
					}

					handle(msg, reader, do)
				}
			}
		}(ctx, kafkaReader, opt)
		readers = append(readers, kafkaReader)
	}
	consumers.store(opt.getAliasName(), readers)
}

func handle(msg kafka.Message, reader *kafka.Reader, do onReceive) {
	innerCtx := context.Background()
	topic := msg.Topic
	actionLog := actionlog.Begin("topic:"+topic, "message-handler")
	actionLog.PutContext("topic", msg.Topic)
	actionLog.PutContext("topic_partition", msg.Partition)
	actionLog.PutContext("topic_offset", msg.Offset)
	key := string(msg.Key)
	if key != "" {
		actionLog.PutContext("key", key)
	}
	var refId, client string
	for _, header := range msg.Headers {
		if header.Key == logKey.RefId {
			refId = string(header.Value)
			continue
		}
		if header.Key == logKey.Client {
			client = string(header.Value)
			continue
		}
	}
	if client != "" {
		actionLog.Client = client
	}
	if refId != "" {
		actionLog.RefId = refId
	}

	innerCtx = context.WithValue(innerCtx, logKey.Id, actionLog.Id)
	innerCtx = context.WithValue(innerCtx, logKey.Action, actionLog.Action)
	actionLog.RequestBody = string(msg.Value)
	logs.DebugWithCtx(innerCtx, "[message] topic: %v, key: %v, message: %v, time: %v, refId: %v, client: %v", topic, string(msg.Key), actionLog.RequestBody, msg.Time, refId, client)

	appkafka.CheckConsumerDelay(innerCtx, msg, actionLog)
	contextMap := make(map[string][]any)
	statMap := make(map[string]float64)
	innerCtx = context.WithValue(innerCtx, logKey.Context, contextMap)
	innerCtx = context.WithValue(innerCtx, logKey.Stat, statMap)

	defer func() {
		if err := recover(); err != nil {
			actionLog.AddStat(statMap)

			actionlog.HandleRecover(err, actionLog, contextMap)
		}
	}()
	do(innerCtx, Message{msg, innerCtx, reader})

	actionLog.AddContext(contextMap)
	actionLog.AddStat(statMap)
	actionlog.End(actionLog, "ok")
}

// Close 关闭消费者
func Close(aliasNames ...string) {
	if len(aliasNames) == 0 {
		aliasNames = append(aliasNames, "default")
	}
	for _, name := range aliasNames {
		kfkConsumers, ok := consumers.load(name)
		if !ok {
			continue
		}
		for _, consumer := range kfkConsumers {
			err := consumer.Close()
			if err != nil {
				log.Printf("kafka consumer <%s> close failed %s", name, err.Error())
			}
		}
		consumers.del(name)
	}
}
