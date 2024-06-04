package logs

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// kafkaLogWriter implements LoggerInterface.
// 将日志写入kafka
type kafkaLogWriter struct {
	producer    *kafka.Writer
	BrokersAddr []string
	Topic       string
	GroupName   string
	Level       LogLevel
}

func newKafkaWriter() Logger {
	w := &kafkaLogWriter{
		Level: LevelDebug,
	}
	return w
}

func (w *kafkaLogWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	w.Topic = opt.Topic
	w.GroupName = opt.GroupName
	w.BrokersAddr = opt.KafkaBrokersAddr
	w.Level = opt.LogLevel
	err := w.startLogger()
	return err
}

func (w *kafkaLogWriter) startLogger() error {
	w.producer = &kafka.Writer{
		Addr:                   kafka.TCP(w.BrokersAddr...),
		Balancer:               new(kafka.LeastBytes),
		RequiredAcks:           kafka.RequireOne,
		Async:                  true,
		AllowAutoTopicCreation: true,
	}
	return nil
}

func (w *kafkaLogWriter) WriteMsg(lm *Msg) error {
	if lm.Level > w.Level {
		return nil
	}

	w.producer.WriteMessages(
		context.TODO(),
		kafka.Message{
			Topic: w.Topic,
			Key:   nil,
			Value: []byte(lm.Format()),
		},
	)
	return nil
}

func (w *kafkaLogWriter) Destroy() {
	_ = w.producer.Close()
}

func init() {
	Register(AdapterKafka, newKafkaWriter)
}
