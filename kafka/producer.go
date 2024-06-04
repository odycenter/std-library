package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	app "std-library/app/conf"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"time"
)

type Producer struct {
	*kafka.Writer
}

// ProducerOption 生产者配置
type ProducerOption struct {
	AliasName         string          `json:"AliasName"`   //别名
	BrokersAddr       []string        `json:"BrokersAddr"` //Kafka单例、集群地址
	Balancer          *kafka.Balancer `json:"-"`           //指定平衡器模式 默认RoundRobin
	Async             bool            `json:"Async"`       //是否异步
	OnDelivery        OnDelivery      `json:"-"`           //同步需要完成该函数，否则线程会阻塞
	SkipTLS           bool            `json:"SkipTLS"`     //跳过TLS验证
	BatchSize         int             `json:"BatchSize"`   //在发送到分区之前限制请求的最大訊息量,默认使用默认值 100。
	BatchTimeout      time.Duration   `json:"-"`           //将不完整的消息批次刷新到 kafka 的时间限制,默认至少每秒刷新一次。
	EnableCompression bool            `json:"-"`
}

func (opt *ProducerOption) getAliasName() string {
	if opt.AliasName == "" {
		opt.AliasName = "default"
	}
	return opt.AliasName
}

func (opt *ProducerOption) getTransport() kafka.RoundTripper {
	if opt.SkipTLS {
		return &kafka.Transport{
			Dial: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).DialContext,
			TLS: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	}
	return kafka.DefaultTransport
}

// WithOnCompletion 设置生产者回调处理
func (opt *ProducerOption) WithOnCompletion(fn OnDelivery) *ProducerOption {
	opt.OnDelivery = fn
	return opt
}

func (opt *ProducerOption) getOnCompletion() onCompletion {
	if opt.OnDelivery != nil {
		return func(messages []kafka.Message, err error) {
			opt.OnDelivery(covertMessagesType(messages), err)
		}
	}
	if !opt.Async && opt.OnDelivery == nil {
		panic(fmt.Errorf("sync producer need to deal the completion message"))
	}
	return nil
}

func (opt *ProducerOption) getBatchSize() int {
	if opt.BatchSize == 0 {
		return 100
	}
	return opt.BatchSize
}

func (opt *ProducerOption) getBatchTimeout() time.Duration {
	if opt.BatchTimeout == 0 {
		return 1 * time.Second
	}
	return opt.BatchTimeout
}

// NewProducer 使用 ProducerOption 创建生产者
func NewProducer(opt *ProducerOption) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(opt.BrokersAddr...),
		Balancer:               new(KeyHashBalancer),
		RequiredAcks:           RequireOne,
		Async:                  opt.Async,
		Completion:             opt.getOnCompletion(),
		Transport:              opt.getTransport(),
		BatchSize:              opt.getBatchSize(),
		BatchTimeout:           opt.getBatchTimeout(),
		AllowAutoTopicCreation: true,
	}
	if opt.EnableCompression {
		writer.Compression = kafka.Snappy
	}
	producers.Store(opt.getAliasName(), &Producer{writer})

	CreateProducer(opt)
}

func covertMessagesType(messages []kafka.Message) (kMessages []DeliveryReport) {
	for _, message := range messages {
		m := DeliveryReport{Topic: message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Key:       message.Key,
			Value:     message.Value,
		}
		kMessages = append(kMessages, m)
	}
	return
}

// Publish 推送数据到目标kafka
// 返回值
// AsyncProducer推送始终返回true
// SyncProducer推送除非推送报错返回false，其他返回true
func (p *Producer) Publish(ctx context.Context, msg MessagePayload) error {
	topic := msg.Topic
	actionLog := actionlog.Begin("topic:"+topic, "message-publisher")
	actionLog.PutContext("topic", msg.Topic)
	key := string(msg.Key)
	if key != "" {
		actionLog.PutContext("key", key)
	}
	message := kafka.Message{
		Topic: topic,
		Key:   msg.Key,
		Value: msg.Value,
		Time:  msg.Timestamp,
	}
	for _, header := range msg.Headers {
		message.Headers = append(message.Headers, kafka.Header{Key: header.Key, Value: header.Value})
	}
	message.Headers = append(message.Headers, kafka.Header{Key: logKey.RefId, Value: []byte(actionLog.Id)})
	if app.Name != "" {
		message.Headers = append(message.Headers, kafka.Header{Key: logKey.Client, Value: []byte(app.Name)})
	}
	actionLog.RequestBody = string(msg.Value)
	defer func() {
		if err := recover(); err != nil {
			actionlog.HandleRecover(err, actionLog, nil)
		}
	}()
	err := p.Writer.WriteMessages(ctx, message)
	if err != nil {
		actionlog.HandleRecover(err, actionLog, nil)
	} else {
		actionlog.End(actionLog, "ok")
	}
	return err
}

func (p *Producer) Close() {
	err := p.Writer.Close()
	if err != nil {
		fmt.Println("producer close error:", err)
	}
}
