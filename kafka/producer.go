package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"time"
)

type Producer struct {
	*kafka.Writer
}

// ProducerOption 生产者配置
type ProducerOption struct {
	AliasName    string          `json:"AliasName"`   //别名
	BrokersAddr  []string        `json:"BrokersAddr"` //Kafka单例、集群地址
	Balancer     *kafka.Balancer `json:"-"`           //指定平衡器模式 默认RoundRobin
	Async        bool            `json:"Async"`       //是否异步
	OnCompletion OnCompletion    `json:"-"`           //同步需要完成该函数，否则线程会阻塞
	SkipTLS      bool            `json:"SkipTLS"`     //跳过TLS验证
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

func (opt *ProducerOption) WithOnCompletion(fn OnCompletion) {
	opt.OnCompletion = fn
}

func (opt *ProducerOption) getOnCompletion() onCompletion {
	if opt.OnCompletion != nil {
		return func(messages []kafka.Message, err error) {
			opt.OnCompletion(covertMessagesType(messages), err)
		}
	}
	if !opt.Async && opt.OnCompletion == nil {
		panic(fmt.Errorf("sync producer need to deal the completion message"))
	}
	return nil
}

// NewProducer 使用 ProducerOption 创建生产者
func NewProducer(opt *ProducerOption) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(opt.BrokersAddr...),
		Balancer:               new(kafka.LeastBytes),
		RequiredAcks:           RequireOne,
		Async:                  opt.Async,
		Completion:             opt.getOnCompletion(),
		Transport:              opt.getTransport(),
		AllowAutoTopicCreation: true,
	}
	producers.Store(opt.getAliasName(), &Producer{writer})
}

// Cli 按 aliasName 获取一个生产者客户端结构题对象
func Cli(aliasName ...string) *Producer {
	name := "default"
	if aliasName != nil {
		name = aliasName[0]
	}
	p, ok := producers.Load(name)
	if !ok {
		log.Panicf("no <%s> kafka client found(need NewProducer first)\n", name)
		return nil
	}
	return p.(*Producer)
}

func covertMessagesType(messages []kafka.Message) (kMessages []Message) {
	for _, message := range messages {
		kMessages = append(kMessages, Message{Message: message})
	}
	return
}

// Send 推送数据到目标kafka
// 返回值
// AsyncProducer推送始终返回true
// SyncProducer推送除非推送报错返回false，其他返回true
func (p *Producer) Send(ctx context.Context, msg kafka.Message) error {
	return p.Writer.WriteMessages(ctx, msg)
}
