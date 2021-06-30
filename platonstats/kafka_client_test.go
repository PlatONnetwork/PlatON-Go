package platonstats

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

//使用方法：
//go test -v -run Test_kafkaClient_producer platonstats/kafka_client_test.go platonstats/confluent_kafka_client.go -args 192.168.9.201:9092 lvxy_test_topic
func Test_kafkaClient_producer(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	fmt.Println("os.Args:", os.Args)
	brokers := os.Args[len(os.Args)-2]
	topic := os.Args[len(os.Args)-1]

	log.Debug("args:", "brokers", brokers, "topic", topic)

	kafkaClient := NewConfluentKafkaClient(brokers, topic)

	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		err := kafkaClient.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0},
			Value:          []byte(word),
		}, nil)
		if err != nil {
			log.Error("cannot queue msg", "err", err)
			fmt.Printf("cannot queue msg, error: %v\n", err)
		}
	}

	// Wait for message deliveries before shutting down
	kafkaClient.producer.Flush(15 * 1000)
}
