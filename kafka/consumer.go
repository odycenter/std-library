package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
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
	AliasName    string    `json:"AliasName"`    //别名
	Count        int       `json:"Count"`        //consumer启动数量
	BrokersAddrs []string  `json:"BrokersAddrs"` //Kafka单例、集群地址
	GroupID      string    `json:"GroupID"`      //消费者组ID
	Topics       []string  `json:"Topics"`       //订阅主题
	MaxBytes     int       `json:"MaxBytes"`     //message最大体积
	ManualCommit bool      `json:"ManualCommit"` //是否手动提交 true手动提交 false自动提交
	Offset       int64     `json:"Offset"`       //起始偏移位
	OnReceive    OnReceive `json:"-"`            //信息消费回调
	OnError      OnError   `json:"-"`            //消费过程中的错误
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

func (opt *ConsumerOption) getCount() int {
	if opt.Count == 0 {
		return 1 // 10MB
	}
	return opt.Count
}

// WithOnReceive 添加消息处理回调
func (opt *ConsumerOption) WithOnReceive(fn OnReceive) {
	opt.OnReceive = fn
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
						return
					}
					do(ctx, Message{msg, ctx, reader})
				}
			}
		}(ctx, kafkaReader, opt)
		readers = append(readers, kafkaReader)
	}
	consumers.store(opt.getAliasName(), readers)
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
