package kafka_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/odycenter/std-library/kafka"
)

func TestProducer(t *testing.T) {
	opt := &kafka.ProducerOption{
		BrokersAddr:  []string{"127.0.0.1:29092"},
		Async:        true,
		OnCompletion: nil,
	}
	opt.WithOnCompletion(func(messages []kafka.Message, err error) {
		if err != nil {
			log.Panicln(err)
		}
		for _, message := range messages {
			log.Printf("send message to %d successful\n", message)
		}
	})
	kafka.NewProducer(opt)
	for range time.Tick(time.Second * 1) {
		err := kafka.Cli().Send(context.Background(), kafka.NewStringMessage("topic1", time.Now().Format(time.RFC3339Nano)))
		if err != nil {
			log.Panicln(err)
		}
	}
}
