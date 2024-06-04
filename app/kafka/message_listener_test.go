package kafka_test

import (
	"context"
	"os"
	app "std-library/app/conf"
	"std-library/app/kafka"
	"testing"
)

type handler struct {
}

func (h *handler) Handle(ctx context.Context, key string, data []byte) {

}

func TestSubscribe(t *testing.T) {
	brokers := []string{"b-1.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092",
		"b-2.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092",
		"b-3.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092"}
	app.Name, _ = os.Hostname()
	kafka.Subscribe(brokers, "topic", &handler{})
}
