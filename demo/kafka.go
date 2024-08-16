package demo

import (
	"context"
	"log/slog"
	actionlog "std-library/app/log"
	"std-library/app/module"
)

type messageHandler struct {
}

func (h *messageHandler) Handle(ctx context.Context, key string, data []byte) {
	slog.Info("key: " + key + ", data: " + string(data))
	//time.Sleep(10 * time.Second)
	actionlog.Context(&ctx, "test", "test123")
	actionlog.Context(&ctx, "test", "test123556666")
}

type KafkaModule struct {
	module.Common
}

func (m *KafkaModule) Initialize() {
	// m.Kafka().GroupId(util.GetIDGenerator().Next(time.Now()))
	m.Kafka().DefaultPoolSize(2)
	m.Kafka().Subscribe("test", &messageHandler{}, 4)
	// a.Load(&KafkaModule{})
	//m.Kafka().Subscribe("kafka-test", &messageHandler{}, 1)

	//kafkaGroup1 := a.Kafka("newOne")
	//kafkaGroup1.GroupId("abcGroup123")
	//kafkaGroup1.Uri(a.RequiredProperty("sys.kafka.uri"))
	//kafkaGroup1.Subscribe("test", &messageHandler{}, 4)

	//m.testKafka()
}

func (m *KafkaModule) testKafka() {
	m.Kafka().Subscribe("player-info-init", &messageHandler{}, 4)
	m.Kafka().Subscribe("pk-combat-log-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-sport-bonus-apply-record-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("create-activity-bonus-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("create-activity-match-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-older-game-transaction", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-daily-task-apply-record-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-lucky-draw-record-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("create-activity-operation-activities-statistic-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("update-activity-operation-activities-statistic-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("create-activity-player-gift-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("update-activity-player-gift-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-sign-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-track-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-track-utility-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-log-statistics-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-notification-popup-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-invite-log-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-sport-order-entity-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-sport-order-virtual-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-vip-level-log-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-activity-sport-apply-record-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-gold-handle-callback-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-refuse-order-callback", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-user-notification-event-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-pay-discount-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-history-data-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-reserve-withdraw-record-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-withdraw-offer-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("steaming-off-live-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("steaming-edit-hot-config-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("cloud-live-activity-player-collection-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("live-activity-player-daily-task", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-notify-live-red-packet-updated-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-live-activity-need-data-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-steaming-hot-log-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("broadcast-live-red-packet-rob-result-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("steamer-player-sub-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("first-pay", &messageHandler{}, 4)
	m.Kafka().Subscribe("receive-field-control-AI-bet-v2", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-steaming-info", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-First-pay-compensation-record", &messageHandler{}, 4)
	m.Kafka().Subscribe("player-invited-gold-log", &messageHandler{}, 4)
	m.Kafka().Subscribe("sync-player-highest-record", &messageHandler{}, 4)
	m.Kafka().Subscribe("receive-live-room-msg-filter-control-AI", &messageHandler{}, 4)
}
