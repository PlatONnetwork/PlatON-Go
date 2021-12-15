package platonstats

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	// historyUpdateRange is the number of blocks a node should report upon login or
	// history request.
	sampleEventChanSize    = 50
	defaultKafkaBlockTopic = "platon-block"
)

type ConfluentKafkaClient struct {
	brokers    string
	blockTopic string
	producer   *kafka.Producer
}

func NewConfluentKafkaClient(urls, blockTopic string) *ConfluentKafkaClient {

	if len(blockTopic) == 0 {
		blockTopic = defaultKafkaBlockTopic
	}

	client := &ConfluentKafkaClient{
		brokers:    urls,
		blockTopic: blockTopic,
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": urls,
		"compression.type":  "gzip",
		"message.max.bytes": 500000000,
	})

	if err != nil {
		log.Error("Failed to create Kafka producer")
		panic(err)
	} else {
		log.Info("Success to create Kafka producer", "urls", urls, "blockTopic", blockTopic)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error("send block message error", "topic", *ev.TopicPartition.Topic, "key", string(ev.Key), "value", string(ev.Value), "err", ev.TopicPartition.Error)
				} else {
					log.Debug("send block message success", "topic", *ev.TopicPartition.Topic, "key", string(ev.Key), "valueSize", len(ev.Value), "value", string(ev.Value))
				}
			}
		}
	}()

	client.producer = p

	fmt.Printf("Success to connect to Kafka")
	return client
}

func (kc *ConfluentKafkaClient) Close() {

	if kc.producer != nil {
		kc.producer.Flush(60 * 1000)
		kc.producer.Close()
	}
}
