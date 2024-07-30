package kafka

import (
	"context"
	kafca "github.com/segmentio/kafka-go"
	"os"
	"runtime/pprof"
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

func BenchmarkKafkaConsumerLog(b *testing.B) {
	msg := kafca.Message{
		Partition: 1,
		Offset:    1,
		Key:       []byte("key"),
		Value:     []byte("{\"ApiName\":\"gametransaction\",\"Data\":{\"LogId\":\"65af212d8cf0a891699b4d8d\",\"PlayerId\":6505538,\"AgentId\":22855,\"ChannelId\":\"C457718_1\",\"PackId\":486,\"TransactionId\":\"65af21271a6d00003b117ccf\",\"SeqNo\":\"65af212d8cf0a891699b4d8c\",\"Type\":11,\"GameId\":8819,\"GameType\":\"SGQP\",\"SubGameId\":8819004,\"AddGold\":23400,\"CreateTime\":1705976109,\"RoundId\":\"24012310144664884863\",\"BetGold\":12000,\"TotalBetGold\":12000,\"WinGold\":23400,\"ValidWater\":12000,\"IsBetTrade\":false,\"Status\":1,\"SettleCount\":1,\"BetTime\":1705976103,\"SettleTime\":1705976103,\"Version\":1,\"KafkaSyncVersion\":1,\"CurrencyCode\":\"RMB\"},\"CreateTime\":\"2024-01-23T10:15:09.589Z\",\"NumberTime\":1705976109,\"Process\":false,\"FailCount\":0,\"AgentId\":0}"),
		Time:      time.Now(),
	}
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.RefId, Value: []byte("8d34199a2342adbbbdf9")})
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.Client, Value: []byte("Client")})
	for i := 0; i < b.N; i++ {
		Handle("client-1", "group1", msg, func(ctx context.Context, key string, data []byte) {
		})
	}
}

func TestKafkaConsumerLog(t *testing.T) {
	msg := kafca.Message{
		Topic:     "topic1234",
		Partition: 1,
		Offset:    1,
		Key:       []byte("key"),
		Value:     []byte("{\"ApiName\":\"gametransaction\",\"Data\":{\"LogId\":\"65af212d8cf0a891699b4d8d\",\"PlayerId\":6505538,\"AgentId\":22855,\"ChannelId\":\"C457718_1\",\"PackId\":486,\"TransactionId\":\"65af21271a6d00003b117ccf\",\"SeqNo\":\"65af212d8cf0a891699b4d8c\",\"Type\":11,\"GameId\":8819,\"GameType\":\"SGQP\",\"SubGameId\":8819004,\"AddGold\":23400,\"CreateTime\":1705976109,\"RoundId\":\"24012310144664884863\",\"BetGold\":12000,\"TotalBetGold\":12000,\"WinGold\":23400,\"ValidWater\":12000,\"IsBetTrade\":false,\"Status\":1,\"SettleCount\":1,\"BetTime\":1705976103,\"SettleTime\":1705976103,\"Version\":1,\"KafkaSyncVersion\":1,\"CurrencyCode\":\"RMB\"},\"CreateTime\":\"2024-01-23T10:15:09.589Z\",\"NumberTime\":1705976109,\"Process\":false,\"FailCount\":0,\"AgentId\":0}"),
		Time:      time.Now(),
	}
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.RefId, Value: []byte("8d34199a2342adbbbdf9")})
	msg.Headers = append(msg.Headers, kafca.Header{Key: logKey.Client, Value: []byte("Client")})

	f, err := os.Create("cpu.pprof")
	if err != nil {
		t.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 10000; i++ {
		Handle("client-1", "group1", msg, func(ctx context.Context, key string, data []byte) {
		})
	}

	f, err = os.Create("mem.pprof")
	if err != nil {
		t.Fatal(err)
	}
	// 写入Memory profiling
	pprof.WriteHeapProfile(f)
	defer f.Close()
}
