package kafka_test

import (
	"context"
	"log"
	"std-library/kafka"
	"testing"
	"time"
)

func TestMessageProducer(t *testing.T) {
	opt := &kafka.ProducerOption{
		BrokersAddr: []string{"b-2.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092",
			"b-3.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092",
			"b-1.mggroupdev.pqz5p8.c2.kafka.ap-northeast-1.amazonaws.com:9092"},
	}
	opt.WithOnCompletion(func(messages []kafka.DeliveryReport, err error) {
		if err != nil {
			log.Panicln(err)
		}
		for _, message := range messages {
			log.Printf("send message to %s successful, partition:%v, key: %s", message.Topic, message.Partition, string(message.Key))
		}
	})

	kafka.CreateProducer(opt)
	payload := "{\"ApiName\":\"loginlog\",\"NumberTime\":1708076109,\"CreateTime\":\"2024-02-16T17:35:09.259763287+08:00\",\"Data\":{\"CreateTime\":1708076109,\"PlayerId\":1194901,\"DeviceId\":\"abc0547ddf8-6196-4a88-898f-3fdcc70d5218\",\"Ip\":\"61.222.239.250\",\"LoginPlatform\":\"uniapp2-web\",\"DeviceModel\":\"pc\",\"SystemVersion\":\"Windows10.0\",\"AppVersion\":\"6.0.15\",\"Logout\":false,\"Reconnect\":true,\"AgentId\":73,\"PackageId\":14,\"ChannelId\":\"dev131231_544\",\"FromDomain\":\"integrate-macau.ljbdev.site\"}}"
	err := kafka.CliV2().Publish(context.Background(), kafka.NewMessageWithKey("sync-login-log", []byte("1194901"), []byte(payload)))
	//payload = "{\"ApiName\":\"gametransaction\",\"Data\":{\"LogId\":\"65cf2dcea9480f0d798aa7fb\",\"PlayerId\":3108079,\"AgentId\":73,\"ChannelId\":\"dev131231_260\",\"PackId\":4,\"TransactionId\":\"65cf2dc7e4350000ad003b0d_73_2258\",\"SeqNo\":\"65cf2dcea9480f0d798aa7fa\",\"Type\":11,\"GameId\":8930,\"GameType\":\"ZBDP\",\"SubGameId\":8930002,\"AddGold\":13650,\"CreateTime\":1708076494,\"RoundId\":\"24021617405406492598\",\"BetGold\":19000,\"TotalBetGold\":19000,\"WinGold\":13650,\"ValidWater\":13000,\"IsBetTrade\":false,\"Status\":1,\"SettleCount\":1,\"BetTime\":1708076487,\"SettleTime\":1708076487,\"Version\":1,\"KafkaSyncVersion\":1,\"CurrencyCode\":\"RMB\",\"UpdateTime\":1708105293},\"CreateTime\":\"2024-02-16T17:41:34.704Z\",\"NumberTime\":1708076494,\"Process\":false,\"FailCount\":0,\"AgentId\":0}"
	//err := kafka.CliV2().Publish(context.Background(), kafka.NewMessageWithKey("gametransaction", []byte("1194901"), []byte(payload)))
	if err != nil {
		log.Panicln(err)
	}

	time.Sleep(8 * time.Second)
}
