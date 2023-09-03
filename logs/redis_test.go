package logs_test

import (
	"log"
	"testing"
	"time"

	"github.com/odycenter/std-library/logs"
)

func TestRedis(t *testing.T) {
	err := logs.SetLogger(logs.AdapterRedis, &logs.Option{
		Adapter:       logs.AdapterRedis,
		LogLevel:      logs.LevelDebug,
		RedisHost:     []string{"127.0.0.1:6379"},
		RedisUsername: "",
		RedisPassword: "",
		IsCluster:     false,
		TLS:           nil,
		RedisKey:      "server_1",
	})
	if err != nil {
		log.Fatalln(err)
	}
	logs.Debug("Debug", 123)
	logs.Info("Info", 123)
	logs.Notice("Notice", 123)
	logs.Warn("Warn", 123)
	logs.Error("Error", 123)
	time.Sleep(time.Second)
}
