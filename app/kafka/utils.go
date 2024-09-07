package kafka

import (
	"context"
	"fmt"
	app "github.com/odycenter/std-library/app/conf"
	actionlog "github.com/odycenter/std-library/app/log"
	"github.com/odycenter/std-library/app/log/consts/logKey"
	"github.com/odycenter/std-library/app/log/dto"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var consumerClientIdSequence atomic.Int32

func Reader(clientID string, opt *SubscribeOption) *kafka.Reader {
	if app.Name == "" { // TODO refactor later
		panic("app.Name is empty, please set it first!")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  clientID,
	}
	kafkaConfig := kafka.ReaderConfig{
		Brokers:     opt.getBrokers(),
		GroupID:     opt.getGroupId(),
		StartOffset: opt.getStartOffset(),
		GroupTopics: []string{opt.getTopic()},
		MaxBytes:    10e5,             // 10MB
		MinBytes:    1,                // 1
		MaxWait:     10 * time.Second, // def:10s
		Dialer:      dialer,
	}
	return kafka.NewReader(kafkaConfig)
}

func ClientID() string {
	podName := os.Getenv("MY_POD_NAME")
	if podName != "" {
		return podName + "-" + strconv.Itoa(int(consumerClientIdSequence.Add(1)))
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host-" + app.Name
	}
	return hostname + "-" + strconv.Itoa(int(consumerClientIdSequence.Add(1)))
}

func CheckConsumerDelay(ctx context.Context, record kafka.Message, actionLog dto.ActionLog) {
	if record.Time.UnixNano() <= 0 {
		return
	}
	consumerDelay := time.Since(record.Time)
	actionLog.PutContext("consumer_delay", consumerDelay.Nanoseconds())
	slog.DebugContext(ctx, fmt.Sprintf("[message] consumer_delay: %v", consumerDelay.String()))
}

func Handle(clientId, groupId string, record kafka.Message, process func(ctx context.Context, key string, data []byte)) (id string) {
	ctx := context.Background()
	topic := record.Topic
	actionLog := actionlog.Begin("topic:"+topic, "message-handler")
	id = actionLog.Id
	if groupId != "" {
		actionLog.PutContext("kafka_group_id", groupId)
	}
	if clientId != "" {
		actionLog.PutContext("kafka_client_id", clientId)
	}
	actionLog.PutContext("topic", record.Topic)
	actionLog.PutContext("topic_partition", record.Partition)
	actionLog.PutContext("topic_offset", record.Offset)

	key := string(record.Key)
	if key != "" {
		actionLog.PutContext("key", key)
	}
	var refId, client, clientHostname string
	for _, header := range record.Headers {
		if header.Key == logKey.RefId {
			refId = string(header.Value)
			continue
		}
		if header.Key == logKey.Client {
			client = string(header.Value)
			continue
		}
		if header.Key == logKey.ClientHostname {
			clientHostname = string(header.Value)
			continue
		}
	}
	contextMap := make(map[string][]any)
	if client != "" {
		actionLog.Client = client
	}
	if refId != "" {
		actionLog.RefId = refId
	}
	if clientHostname != "" {
		actionLog.PutContext(logKey.ClientHostname, clientHostname)
	}
	statMap := make(map[string]float64)
	ctx = context.WithValue(ctx, logKey.Stat, statMap)
	ctx = context.WithValue(ctx, logKey.Context, contextMap)
	ctx = context.WithValue(ctx, logKey.Id, id)
	ctx = context.WithValue(ctx, logKey.Action, actionLog.Action)
	actionLog.RequestBody = string(record.Value)
	slog.DebugContext(ctx, fmt.Sprintf("[message] topic: %v, key: %v, message: %v, time: %v, refId: %v, client: %v", topic, string(record.Key), actionLog.RequestBody, record.Time, refId, client))

	CheckConsumerDelay(ctx, record, actionLog)

	defer func() {
		if err := recover(); err != nil {
			actionLog.AddStat(statMap)

			actionlog.HandleRecover(err, actionLog, contextMap)
		}
	}()

	process(ctx, key, record.Value)
	actionLog.AddContext(contextMap)
	actionLog.AddStat(statMap)
	actionlog.End(actionLog, "ok")

	return
}
