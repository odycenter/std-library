package logs_test

import (
	"log"
	"testing"
	"time"

	"github.com/odycenter/std-library/logs"
)

func TestKafka(t *testing.T) {
	logs.NewLogger(10000)
	err := logs.SetLogger(logs.AdapterKafka, &logs.Option{
		KafkaBrokersAddr: []string{"localhost:29092"},
		Topic:            "server_name_1",
		Adapter:          logs.AdapterKafka,
		LogLevel:         logs.LevelDebug,
	})
	if err != nil {
		log.Fatalln(err)
	}
	logs.Info("asddsdasd")
	logs.Error("adsadsdas")
	time.Sleep(time.Second * 5)
}
