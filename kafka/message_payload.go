package kafka

import "time"

type Header struct {
	Key   string
	Value []byte
}

type MessagePayload struct {
	Topic     string
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Headers   []Header
}

func NewMessage(topic string, value []byte) MessagePayload {
	return MessagePayload{
		Topic: topic,
		Value: value,
	}
}

func NewMessageWithKey(topic string, key, value []byte) MessagePayload {
	return MessagePayload{
		Topic: topic,
		Key:   key,
		Value: value,
	}
}

// NewStringMessage 创建新的string Message
func NewStringMessage(topic string, value string) MessagePayload {
	return MessagePayload{
		Topic: topic,
		Value: []byte(value),
	}
}
