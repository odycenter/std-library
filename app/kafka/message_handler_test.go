package kafka

import (
	"context"
	kafca "github.com/segmentio/kafka-go"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/web/errors"
	"std-library/logs"
	"testing"
	"time"
)

func TestHandle(t *testing.T) {
	msg := kafca.Message{
		Partition: 1,
		Offset:    1,
		Topic:     "topic",
		Key:       []byte("key"),
		Value:     []byte("{\"ApiName\":\"gametransaction\"}"),
		Time:      time.Now(),
	}
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.RefId, Value: []byte("8d34199a2342adbbbdf9")})
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.Client, Value: []byte("Client")})
	Handle("client-1", "group1", msg, func(ctx context.Context, key string, data []byte) {
		logs.InfoWithCtx(ctx, "key: %s, data: %s", key, data)
		actionlog.Context(&ctx, "topic123", "msg.Topic")
		errors.InternalError(-1, "error")
	})
}
