package kafka

import "context"

type MessageHandler interface {
	Handle(ctx context.Context, key string, data []byte)
}
