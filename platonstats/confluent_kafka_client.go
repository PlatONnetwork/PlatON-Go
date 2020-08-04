package platonstats

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	// historyUpdateRange is the number of blocks a node should report upon login or
	// history request.
	sampleEventChanSize                      = 50
	defaultKafkaBlockTopic                   = "platon-block"
	defaultKafkaAccountCheckingConsumerGroup = "platon-account-checking-group"
	defaultKafkaAccountCheckingTopic         = "platon-account-checking"
)

type ConfluentKafkaClient struct {
	brokers                      string
	blockTopic                   string
	accountCheckingTopic         string
	AccountCheckingConsumerGroup string
	producer                     *kafka.Producer
	consumer                     *kafka.Consumer
}

func NewConfluentKafkaClient(urls, blockTopic, checkingTopic, checkingConsumerGroup string) *ConfluentKafkaClient {
	//brokers := strings.Split(urls, ",")

	if len(blockTopic) == 0 {
		blockTopic = defaultKafkaAccountCheckingTopic
	}

	if len(checkingTopic) == 0 {
		checkingTopic = defaultKafkaAccountCheckingTopic
	}

	if len(checkingConsumerGroup) == 0 {
		checkingConsumerGroup = defaultKafkaAccountCheckingConsumerGroup
	}

	client := &ConfluentKafkaClient{
		brokers:                      urls,
		blockTopic:                   blockTopic,
		accountCheckingTopic:         checkingTopic,
		AccountCheckingConsumerGroup: checkingConsumerGroup,
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": urls})

	if err != nil {
		log.Error("Failed to create Kafka producer")
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error("send block message error", "key", ev.Key, "value", ev.Value, "err", ev.TopicPartition.Error)
				} else {
					log.Debug("send block message success", "key", ev.Key, "value", ev.Value)
				}
			}
		}
	}()

	client.producer = p

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": urls,
		"group.id":          checkingConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Error("Failed to create Kafka consumer")
		panic(err)
	}

	err = c.Subscribe(checkingTopic, nil)
	if err != nil {
		log.Error("Failed to subscribe consumer topic")
		panic(err)
	}

	client.consumer = c
	fmt.Printf("Success to connect to Kafka")
	return client
}

func (kc *ConfluentKafkaClient) Close() {

	if kc.producer != nil {
		kc.producer.Close()
	}
	if kc.consumer != nil {
		if err := kc.consumer.Close(); err != nil {
			log.Error("Failed to close consumer", "err", err)
		}
	}
}
