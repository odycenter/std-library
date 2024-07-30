package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type OnCompletion func(messages []Message, err error)

type onCompletion func(messages []kafka.Message, err error)

type OnError func(err error)

type Message struct {
	kafka.Message
	ctx    context.Context
	reader *kafka.Reader
}
