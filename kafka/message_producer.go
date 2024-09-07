package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/odycenter/std-library/app/async"
	app "github.com/odycenter/std-library/app/conf"
	actionlog "github.com/odycenter/std-library/app/log"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	"log/slog"
	"strings"
)

type MessageProducer interface {
	Publish(ctx context.Context, msg MessagePayload) error
	Close()
}

type messageProducer struct {
	producer *kafka.Producer
}

var DisableForceProducerFlush = false

func (p *messageProducer) Publish(ctx context.Context, msg MessagePayload) error {
	topic := msg.Topic
	actionLog := actionlog.Begin("topic:"+topic, "message-publisher")
	actionName := "topic:" + topic
	if ctx != nil {
		rootAction := actionlog.GetAction(&ctx)
		if rootAction != "" {
			actionName = rootAction + ":" + actionName
			actionLog.PutContext("root_action", rootAction)
			actionLog.RefId = actionlog.GetId(&ctx)
		}
	}
	actionLog.Action = actionName
	actionLog.PutContext("topic", msg.Topic)
	actionLog.PutContext("kafka_client_version", "v2")
	key := string(msg.Key)
	if key != "" {
		actionLog.PutContext("key", key)
	}
	headers := make([]kafka.Header, 0, len(msg.Headers))
	for _, header := range msg.Headers {
		headers = append(headers, kafka.Header{Key: header.Key, Value: header.Value})
	}
	headers = append(headers, kafka.Header{Key: logKey.RefId, Value: []byte(actionLog.Id)})
	headers = append(headers, kafka.Header{Key: logKey.ClientHostname, Value: []byte(app.LocalHostName())})
	if app.Name != "" {
		headers = append(headers, kafka.Header{Key: logKey.Client, Value: []byte(app.Name)})
	}
	actionLog.RequestBody = string(msg.Value)
	defer func() {
		if err := recover(); err != nil {
			actionlog.HandleRecover(err, actionLog, nil)
		}
	}()

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            msg.Key,
		Value:          msg.Value,
		Headers:        headers,
	}, nil)

	if err != nil {
		actionlog.HandleRecover(err, actionLog, nil)
	} else {
		actionlog.End(actionLog, "ok")
	}
	if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrQueueFull {
		if !DisableForceProducerFlush {
			async.RunFuncWithName(&ctx, "force-flush-message", func(ctx context.Context) {
				slog.ErrorContext(ctx, "Kafka local queue full error!! trying to Flush ...")
				flushedMessages := p.producer.Flush(30 * 1000)
				slog.InfoContext(ctx, fmt.Sprintf("Flushed kafka messages. Outstanding events still un-flushed: %d", flushedMessages))
			})
		}
	}
	return err
}

func (p *messageProducer) CreateTopic(topic string, numPartitions int) error {
	client, err := kafka.NewAdminClientFromProducer(p.producer)
	if err != nil {
		return err
	}

	topics, err := client.CreateTopics(context.Background(), []kafka.TopicSpecification{{
		Topic:         topic,
		NumPartitions: numPartitions,
	}})
	if err != nil {
		return err
	}
	for _, topic := range topics {
		if topic.Error.Code() != kafka.ErrNoError {
			return topic.Error
		}
	}
	return nil
}

func (p *messageProducer) Close() {
	p.Close()
}

func CreateProducer(opt *ProducerOption) {
	servers := strings.Join(opt.BrokersAddr, ",")
	config := &kafka.ConfigMap{
		"bootstrap.servers": servers,
	}
	if opt.EnableCompression {
		err := config.SetKey("compression.type", "snappy")
		if err != nil {
			panic(err)
		}
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	if opt.OnDelivery != nil {
		go func() {
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					m := ev
					message := DeliveryReport{
						Topic:     *m.TopicPartition.Topic,
						Partition: int(m.TopicPartition.Partition),
						Offset:    int64(m.TopicPartition.Offset),
						Key:       m.Key,
						Value:     m.Value,
					}
					if m.TopicPartition.Error != nil {
						slog.Warn(fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error))
					}
					opt.OnDelivery([]DeliveryReport{message}, m.TopicPartition.Error)
				default:
					continue
				}
			}
		}()
	}

	v2Producers.Store(opt.getAliasName(), &messageProducer{producer: p})
}

type DeliveryReport struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
}

type OnDelivery func(messages []DeliveryReport, err error)
