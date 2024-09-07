package kafka_test

import (
	"context"
	"github.com/odycenter/std-library/kafka"
	"log"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	opt := &kafka.ProducerOption{
		BrokersAddr: []string{"127.0.0.1:29092"},
		Async:       true,
		OnDelivery:  nil,
	}
	opt.WithOnCompletion(func(messages []kafka.DeliveryReport, err error) {
		if err != nil {
			log.Panicln(err)
		}
		for _, message := range messages {
			log.Printf("send message to %s successful\n", message.Topic)
		}
	})
	kafka.NewProducer(opt)
	for range time.Tick(time.Second * 1) {
		err := kafka.Cli().Publish(context.Background(), kafka.NewStringMessage("topic1", time.Now().Format(time.RFC3339Nano)))
		if err != nil {
			log.Panicln(err)
		}
	}
}
