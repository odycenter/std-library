package kafka

import (
	"github.com/segmentio/kafka-go"
	"sync"
)

var producers sync.Map     //map[string]UniversalProducer
var consumers consumersMap //map[string][]Consumer

// 起始偏移位置
const (
	// FirstOffset The most recent offset available for a partition.
	FirstOffset = kafka.FirstOffset

	// LastOffset The least recent offset available for a partition.
	LastOffset = kafka.LastOffset
)

// 是否等待ACK响应
const (
	RequireNone = kafka.RequireNone //(0) 不等待确认 默认
	RequireOne  = kafka.RequireOne  //(1) 等待leader确认
	RequireAll  = kafka.RequireAll  //(-1) 等待全体确认
)

// balancer分区分配均衡器
var (
	//RoundRobin 是一种 Balancer 实现，它在所有可用分区之间平均分配消息。
	RoundRobin = new(kafka.RoundRobin)
	//LeastBytes 是一个Balancer实现，它将消息路由到接收到最少数据量的分区。
	//请注意，多个生产者之间没有进行协调，良好的平衡依赖于每个使用 LeastBytes 平衡器的生产者都应该产生良好平衡的消息这一事实。
	//(*LeastBytes) 上的方法：
	//Balance(msg Message, partitions ...int) int
	//makeCounters(partitions ...int) (counters []leastBytesCounter)
	LeastBytes = new(kafka.LeastBytes)
	// Hash 它使用提供的哈希函数来确定将消息路由到哪个分区。
	//这确保具有相同密钥的消息被路由到相同的分区。
	//计算分区的逻辑是：
	//hasher.Sum32() % len(partitions) => partition
	//Hash默认使用FNV-1a算法。
	//这与 Sarama Producer 使用的算法相同，并确保 kafka-go 生成的消息将被传送到 Sarama producer 将被传送到的相同主题。
	Hash = new(kafka.Hash)
	// ReferenceHash 使用提供的哈希函数来确定将消息路由到哪个分区。这确保具有相同密钥的消息被路由到相同的分区。
	//计算分区的逻辑是：
	//(int32(hasher.Sum32()) & 0x7fffffff) % len(partitions) => partition
	//默认情况下，ReferenceHash 使用 FNV-1a 算法。这是与 Sarama NewReferenceHashPartitioner 相同的算法，并确保由 kafka-go 生成的消息将被传送到与 Sarama 生产者将被传送到的相同主题。
	//(*ReferenceHash) 上的方法：
	//Balance(msg Message, partitions ...int) int
	ReferenceHash = new(kafka.ReferenceHash)
	// CRC32Balancer 使用 CRC32 哈希函数来确定将消息路由到哪个分区。这确保具有相同密钥的消息被路由到相同的分区。
	//这个平衡器与 librdkafka 中的内置哈希分区器以及构建在它之上的语言绑定兼容，包括 github.com/confluentinc/confluent-kafka-go Go 包。
	//Consistent字段为 false（默认）时，此分区程序等效于 librdkafka 中的“consistent_random”设置。当Consistent为真时，此分区程序等效于“一致”设置。后者会将空键或零键散列到同一分区中。
	//除非您绝对确定您的所有消息都有密钥，否则最好关闭Consistent标志。否则，您可能会创建一个非常热的分区。
	//关于（CRC32Balancer）的方法：
	//Balance(msg Message, partitions ...int) (partition int)
	CRC32Balancer = new(kafka.CRC32Balancer)
	//Murmur2Balancer 使用 Murmur2 哈希函数来确定将消息路由到哪个分区。这确保具有相同密钥的消息被路由到相同的分区。此平衡器与 Java 库和 librdkafka 的“murmur2”和“murmur2_random”分区器使用的分区器兼容。
	//Consistent字段为 false（默认）时，此分区程序等效于 librdkafka 中的“murmur2_random”设置。当Consistent为真时，此分区程序等效于“murmur2”设置。后者会将 nil 键散列到同一个分区中。无论配置如何，空的非零键总是散列到同一个分区。
	//除非您绝对确定您的所有消息都有密钥，否则最好关闭Consistent标志。否则，您可能会创建一个非常热的分区。
	//请注意，librdkafka 文档指出“murmur2_random”在功能上等同于默认的 Java 分区程序。那是因为 Java 分区器将使用循环平衡器而不是对 nil 键进行随机分配。我们选择 librdkafka 的实现，因为它可以说具有更大的安装基础。
	Murmur2Balancer = new(kafka.Murmur2Balancer)
)
