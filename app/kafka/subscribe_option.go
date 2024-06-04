package kafka

import "github.com/segmentio/kafka-go"

type SubscribeOption struct {
	Brokers     []string
	GroupId     string
	Topic       string
	StartOffset int64
}

func (opt *SubscribeOption) getStartOffset() int64 {
	if opt.StartOffset == 0 {
		return kafka.FirstOffset
	}
	return opt.StartOffset
}

func (opt *SubscribeOption) getGroupId() string {
	if opt.GroupId == "" {
		panic("GroupID is empty")
	}
	return opt.GroupId
}

func (opt *SubscribeOption) getTopic() string {
	if opt.Topic == "" {
		panic("topic is empty")
	}
	return opt.Topic
}

func (opt *SubscribeOption) getBrokers() []string {
	if opt.Brokers == nil || len(opt.Brokers) == 0 {
		panic("brokers is empty")
	}
	return opt.Brokers
}
