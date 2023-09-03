package kafka_test

import (
	"context"
	"log"
	"std-library/kafka"
	"testing"
	"time"
)

func TestKafkaConsumer(t *testing.T) {
	kafkaConfig := kafka.ConsumerOption{
		Count:        3,
		BrokersAddrs: []string{"127.0.0.1:29092"},
		GroupID:      "test1",
		Topics:       []string{"topic1"},
		MaxBytes:     10e6, // 10MB
		ManualCommit: false,
		OnReceive: func(ctx context.Context, msg kafka.Message) {
			log.Printf("key:%s,value:%s,offset:%d,position:%d\n", msg.Key, msg.Value, msg.Offset, msg.Partition)
			//time.Sleep(time.Millisecond * 1)
		},
		OnError: func(err error) {
			log.Println("receive the error from Kafka", err.Error())
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kafka.NewConsumer(ctx, &kafkaConfig)
	<-time.After(time.Hour)
}
