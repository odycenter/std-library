package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	internal "std-library/app/internal/module"
	"std-library/logs"
	"sync"
	"sync/atomic"
	"time"
)

type MessageListener struct {
	topic        string
	Opt          *SubscribeOption
	Handler      MessageHandler
	ctx          context.Context
	cancel       context.CancelFunc
	reader       []*kafka.Reader
	poolSize     int
	runningTasks int32
	mu           sync.Mutex
}

func (m *MessageListener) SetPoolSize(size int) {
	m.poolSize = size
}

func (m *MessageListener) PoolSize() int {
	return m.poolSize
}

func (m *MessageListener) Initialize(opt *SubscribeOption) {
	m.Opt = opt
	m.topic = opt.Topic
	m.ctx, m.cancel = context.WithCancel(context.Background())
}

func (m *MessageListener) Start(ctx context.Context) {
	for i := 0; i < m.poolSize; i++ {
		clientId := ClientID()
		logs.InfoWithCtx(ctx, "[message-listener] start message listener, groupId: %s, topic: %s, clientID: %s", m.Opt.GroupId, m.topic, clientId)
		go m.Run(clientId)
	}
}

func (m *MessageListener) Run(clientId string) {
	reader := Reader(clientId, m.Opt)
	m.mu.Lock()
	m.reader = append(m.reader, reader)
	m.mu.Unlock()
	for {
		select {
		case <-m.ctx.Done():
			logs.Warn("[message-listener] ctx.Done by caller. groupId: %s, topic: %s, clientID: %s", m.Opt.GroupId, m.topic, clientId)
			return
		default:
			if internal.IsShutdown() {
				logs.Info("[message-listener] reject kafka handle process due to server is shutting down!! GroupId: %s, topic: %s, clientId: %s", m.Opt.GroupId, m.topic, clientId)
				return
			}
			m.run(clientId, reader)
		}
	}
}

func (m *MessageListener) run(clientId string, reader *kafka.Reader) {
	atomic.AddInt32(&m.runningTasks, 1)
	defer atomic.AddInt32(&m.runningTasks, -1)

	innerCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	msg, err := reader.FetchMessage(innerCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return
		}
		logs.Error("[message-listener] FetchMessage fail, groupId: %s, topic: %s, clientId: %s, error:", m.Opt.GroupId, m.topic, clientId, err)
		return
	}

	id := Handle(clientId, m.Opt.GroupId, msg, m.Handler.Handle)
	err = reader.CommitMessages(context.Background(), msg)
	if err != nil {
		logs.Error("[message-listener] CommitMessages fail, groupId: %s, topic: %s, clientId: %s, id: %v, error: ", m.Opt.GroupId, m.topic, clientId, id, err)
	}
}

func (m *MessageListener) RunningTasks() int {
	return int(atomic.LoadInt32(&m.runningTasks))
}

func (m *MessageListener) AwaitTermination(ctx context.Context, timeoutInMs int64) {
	logs.InfoWithCtx(ctx, "shutting down message listener. groupId: %s, topic: %s", m.Opt.GroupId, m.topic)

	innerCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInMs)*time.Millisecond)
	defer cancel()

	for {
		select {
		case <-innerCtx.Done():
			logs.InfoWithCtx(innerCtx, "[FAILED_TO_STOP] failed to terminate message listener, due to timeout, groupId: %s, topic: %s, canceledTasks=%d", m.Opt.GroupId, m.topic, m.RunningTasks())
			m.cancel()
			for _, reader := range m.reader {
				go reader.Close()
			}
			return
		default:
			if m.RunningTasks() == 0 {
				logs.InfoWithCtx(innerCtx, "all message handler have completed, groupId: %s, topic: %s", m.Opt.GroupId, m.topic)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

}
