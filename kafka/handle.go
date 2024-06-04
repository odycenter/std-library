package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type OnReceive func(ctx context.Context, msg Message)

type onReceive func(ctx context.Context, msg Message)

type OnCompletion func(messages []Message, err error)

type onCompletion func(messages []kafka.Message, err error)

type OnError func(err error)

type Message struct {
	kafka.Message
	ctx    context.Context
	reader *kafka.Reader
}

// Commit 确认已被消费
func (msg *Message) Commit() error {
	return msg.reader.CommitMessages(msg.ctx, msg.Message)
}
