package logs

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// kafkaLogWriter implements LoggerInterface.
// 将日志写入kafka
type kafkaLogWriter struct {
	producer     *kafka.Writer
	BrokersAddr  []string
	Topic        string
	GroupName    string
	Level        LogLevel
	Formatter    string
	logFormatter LogFormatter
}

func newKafkaWriter() Logger {
	w := &kafkaLogWriter{
		Level: LevelDebug,
	}
	w.logFormatter = w
	return w
}

func (w *kafkaLogWriter) Format(lm *Msg) string {
	return lm.Format()
}

func (w *kafkaLogWriter) SetFormatter(f LogFormatter) {
	w.logFormatter = f
}

func (w *kafkaLogWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	w.Topic = opt.Topic
	w.GroupName = opt.GroupName
	w.BrokersAddr = opt.KafkaBrokersAddr
	w.Level = opt.LogLevel
	w.Formatter = opt.Formatter
	if len(w.Formatter) > 0 {
		fmtr, ok := GetFormatter(w.Formatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", w.Formatter)
		}
		w.logFormatter = fmtr
	}
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

	msg := w.logFormatter.Format(lm)
	w.producer.WriteMessages(
		context.TODO(),
		kafka.Message{
			Topic: w.Topic,
			Key:   nil,
			Value: []byte(msg),
		},
	)
	return nil
}

func (w *kafkaLogWriter) Destroy() {
	_ = w.producer.Close()
}

func (w *kafkaLogWriter) Flush() {

}

func init() {
	Register(AdapterKafka, newKafkaWriter)
}
