package kafka

type ConsumerOption struct {
	BrokersAddrs []string `json:"BrokersAddrs"` //Kafka单例、集群地址
	GroupID      string   `json:"GroupID"`      //消费者组ID
}
